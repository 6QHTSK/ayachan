package Services

import (
	"fmt"
	"github.com/6QHTSK/ayachan/Config"
	"github.com/6QHTSK/ayachan/Databases"
	"github.com/6QHTSK/ayachan/Log"
	"github.com/6QHTSK/ayachan/Models"
	"github.com/6QHTSK/ayachan/Models/ChartFormat"
	"github.com/6QHTSK/ayachan/utils"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type chartDataRequest struct {
	Result bool                          `json:"result"`
	Info   ChartFormat.BestdoriChartItem `json:"info"`
}

type chartDataList struct {
	Result bool `json:"result"`
	Count  int  `json:"count"`
	List   []struct {
		ChartID int                `json:"id"`
		Author  ChartFormat.Author `json:"author"`
		Diff    int                `json:"diff"`
		Level   int                `json:"level"`
		Likes   int                `json:"likes"`
	} `json:"list"`
}

func BestdoriFanMadeSyncAll() (errCode int, err error) {
	totalCount, errCode, err := BestdoriFanMadeSyncPage(0)
	if err != nil {
		return errCode, err
	}
	totalPage := int(math.Ceil(float64(totalCount) / 50.0))
	for i := 1; i < totalPage; i++ {
		Log.Log.Debugf("Page %d/%d\n", i, totalPage)
		_, errCode, err := BestdoriFanMadeSyncPage(i)
		if err != nil {
			return errCode, err
		}
	}
	Log.Log.Info("BestdoriFanMadeSyncAll Success")
	return http.StatusOK, nil
}

func BestdoriFanMadeSyncRand() (errCode int, err error) {
	totalCount, errCode, err := BestdoriFanMadeSyncPage(0)
	if err != nil {
		return errCode, err
	}
	totalPage := int(math.Ceil(float64(totalCount) / 50.0))
	syncList := []int{1, 2, 3, rand.Intn(totalPage), rand.Intn(totalPage), rand.Intn(totalPage)}
	for _, item := range syncList {
		Log.Log.Debugf("Sync Page %d", item)
		_, errCode, err := BestdoriFanMadeSyncPage(item)
		if err != nil {
			return errCode, err
		}
	}
	Log.Log.Info("BestdoriFanMadeSyncRand Success")
	return http.StatusOK, nil
}

func BestdoriFanMadeSyncPage(page int) (totalCount int, errCode int, err error) {
	listDataParam, err := url.Parse(fmt.Sprintf("list?page=%d&limit=50", page))
	listDataUrl := Config.BestdoriAPIUrl.ResolveReference(listDataParam)
	var request chartDataList
	for i := 1; i <= 5; i++ {
		errCode, err = utils.HttpGet(listDataUrl.String(), &request)
		if err == nil {
			break
		}
		if err != nil {
			Log.Log.Warningf("Failed to fetch page info %d [Attempt %d]", page, i)
			if i == 5 {
				return totalCount, errCode, err
			}
		}
	}

	var wg sync.WaitGroup
	ch := make(chan bool, 7)
	for i, item := range request.List {
		res, err := Databases.CheckBestdoriSongVersion(item.ChartID)
		if err != nil {
			return totalCount, http.StatusInternalServerError, err
		}
		if res {
			// Update Author's nickname & like count
			err = Databases.UpdateBestdori(ChartFormat.BestdoriChartUpdateItem{
				ChartID:  item.ChartID,
				Username: item.Author.Username,
				Nickname: item.Author.Nickname,
				Diff:     item.Diff,
				Level:    item.Level,
				Likes:    item.Likes,
			})
			if err != nil {
				return totalCount, http.StatusInternalServerError, err
			}
		} else {
			go func(i int, item int, ch chan bool) {
				defer func() {
					<-ch
					err := recover()
					if err != nil {
						Log.Log.Errorf("Panic While Updating Chart #%d : %s", item, err)
					}
					wg.Done()
				}()
				ch <- true
				wg.Add(1)
				var j int
				for j = 1; j <= 5; j++ {
					errCode, err = BestdoriFanMadeInsertID(item)
					if err == nil {
						Log.Log.Tracef("Success Update Chart %d [Attempt %d]", item, j)
						return
					} else {
						Log.Log.Warningf("Failed to update Chart %d [Attempt %d] : Error %s", item, j, err.Error())
					}
				}
				Log.Log.Warningf("Attempt times exceed!")
			}(i, item.ChartID, ch)
		}
	}
	time.Sleep(time.Second * 2)
	wg.Wait()
	return request.Count, http.StatusOK, nil
}

func BestdoriFanMadeInsertID(chartID int) (errorCode int, err error) {
	chartDataParam, err := url.Parse(fmt.Sprintf("%d", chartID))
	chartDataUrl := Config.BestdoriAPIUrl.ResolveReference(chartDataParam)
	var request chartDataRequest
	errorCode, err = utils.HttpGet(chartDataUrl.String(), &request)
	if err != nil {
		return errorCode, err
	}
	bestdoriChartItem := request.Info

	result, err := request.Info.Chart.MapCheck()
	if !result {
		Log.Log.Warningf("谱面无法解析,%s", err)
		return http.StatusBadRequest, err
	}
	Map := request.Info.Chart.Decode()

	// Insert
	bestdoriChartItem.MapMetricsBasic, _, _, _, _ = basicMetricsGetter(Map)
	_, bestdoriChartItem.IrregularInfo = ParseMap(Map)

	if bestdoriChartItem.IrregularInfo.Irregular == Models.RegularTypeUnknown {
		Log.Log.Errorf("分析异常chartID：%d", chartID)
	}

	RuneContent := []rune(bestdoriChartItem.Content)
	if len(RuneContent) > 800 {
		bestdoriChartItem.Content = string([]rune(bestdoriChartItem.Content)[:800])
	}
	err = Databases.InsertBestdori(bestdoriChartItem)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// MysqlSyncToMeiliSearch call at 30s each minute
func MysqlSyncToMeiliSearch() (err error) {
	lastUpdate, err := Databases.GetMeiliLastUpdate()
	if err != nil {
		return err
	}
	documents, err := Databases.QueryBestdoriFanMadeByLastUpdate(lastUpdate)
	if err != nil {
		return err
	}
	if len(documents) > 0 {
		count := len(documents)
		for i := 0; i < len(documents); i += 500 {
			end := i + 500
			if i+500 > len(documents) {
				end = len(documents)
			}
			err := Databases.AddDocument(documents[i:end])
			if err != nil {
				for _, doc := range documents {
					err := Databases.AddDocument(doc)
					if err != nil {
						Log.Log.Warningf("Doc[%d] Add Failed, err: %s", doc.ChartID, err)
						count--
					}
				}
			}
		}
		Log.Log.Infof("MeiliSearch Sync Finish! %d documents sync!", count)
	} else {
		// Nothing to sync
		Log.Log.Info("Nothing to do, MeiliSearch Sync Sleep")
	}
	return nil
}

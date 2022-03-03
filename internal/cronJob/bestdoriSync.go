package cronJob

import (
	"fmt"
	"github.com/6QHTSK/ayachan/internal/cronJob/config"
	"github.com/6QHTSK/ayachan/internal/pkg/generalModels"
	"github.com/6QHTSK/ayachan/internal/pkg/httpx"
	"github.com/6QHTSK/ayachan/internal/pkg/logrus"
	"github.com/6QHTSK/ayachan/pkg"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type chartDataRequest struct {
	Result bool                            `json:"result"`
	Info   generalModels.BestdoriChartItem `json:"info"`
}

type chartDataList struct {
	Result bool `json:"result"`
	Count  int  `json:"count"`
	List   []struct {
		ChartID int                  `json:"id"`
		Author  generalModels.Author `json:"author"`
		Diff    int                  `json:"diff"`
		Level   int                  `json:"level"`
		Likes   int                  `json:"likes"`
	} `json:"list"`
}

func BestdoriFanMadeSyncAll() (errCode int, err error) {
	totalCount, errCode, err := BestdoriFanMadeSyncPage(0)
	if err != nil {
		return errCode, err
	}
	totalPage := int(math.Ceil(float64(totalCount) / 50.0))
	for i := 1; i < totalPage; i++ {
		logrus.Log.Debugf("Page %d/%d\n", i, totalPage)
		_, errCode, err := BestdoriFanMadeSyncPage(i)
		if err != nil {
			return errCode, err
		}
	}
	logrus.Log.Info("BestdoriFanMadeSyncAll Success")
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
		logrus.Log.Debugf("Sync Page %d", item)
		_, errCode, err := BestdoriFanMadeSyncPage(item)
		if err != nil {
			return errCode, err
		}
	}
	logrus.Log.Info("BestdoriFanMadeSyncRand Success")
	return http.StatusOK, nil
}

func BestdoriFanMadeSyncPage(page int) (totalCount int, errCode int, err error) {
	listDataParam, err := url.Parse(fmt.Sprintf("list?page=%d&limit=50", page))
	listDataUrl := config.BestdoriAPIUrl.ResolveReference(listDataParam)
	var request chartDataList
	for i := 1; i <= 5; i++ {
		errCode, err = httpx.HttpGet(listDataUrl.String(), &request)
		if err == nil {
			break
		}
		if err != nil {
			logrus.Log.Warningf("Failed to fetch page info %d [Attempt %d]", page, i)
			if i == 5 {
				return totalCount, errCode, err
			}
		}
	}

	var wg sync.WaitGroup
	ch := make(chan bool, 7)
	for i, item := range request.List {
		res, err := CheckBestdoriSongVersion(item.ChartID)
		if err != nil {
			return totalCount, http.StatusInternalServerError, err
		}
		if res {
			// Update Author's nickname & like count
			err = UpdateBestdori(generalModels.BestdoriChartUpdateItem{
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
						logrus.Log.Errorf("Panic While Updating Chart #%d : %s", item, err)
					}
					wg.Done()
				}()
				ch <- true
				wg.Add(1)
				var j int
				for j = 1; j <= 5; j++ {
					errCode, err = BestdoriFanMadeInsertID(item)
					if err == nil {
						logrus.Log.Tracef("Success Update Chart %d [Attempt %d]", item, j)
						return
					} else {
						logrus.Log.Warningf("Failed to update Chart %d [Attempt %d] : Error %s", item, j, err.Error())
					}
				}
				logrus.Log.Warningf("Attempt times exceed!")
			}(i, item.ChartID, ch)
		}
	}
	time.Sleep(time.Second * 2)
	wg.Wait()
	return request.Count, http.StatusOK, nil
}

func BestdoriFanMadeInsertID(chartID int) (errorCode int, err error) {
	chartDataParam, err := url.Parse(fmt.Sprintf("%d", chartID))
	chartDataUrl := config.BestdoriAPIUrl.ResolveReference(chartDataParam)
	var request chartDataRequest
	errorCode, err = httpx.HttpGet(chartDataUrl.String(), &request)
	if err != nil {
		return errorCode, err
	}
	bestdoriChartItem := request.Info

	result, err := request.Info.Chart.MapCheck()
	if !result {
		logrus.Log.Warningf("谱面无法解析,%s", err)
		return http.StatusBadRequest, err
	}
	Map := request.Info.Chart.Decode()

	// Insert
	bestdoriChartItem.MapMetricsBasic = pkg.StandardInfoGetter(Map).MapMetricsBasic
	_, bestdoriChartItem.IrregularInfo = pkg.ParseMap(Map)

	if bestdoriChartItem.IrregularInfo.Irregular == pkg.RegularTypeUnknown {
		logrus.Log.Errorf("分析异常chartID：%d", chartID)
	}

	RuneContent := []rune(bestdoriChartItem.Content)
	if len(RuneContent) > 800 {
		bestdoriChartItem.Content = string([]rune(bestdoriChartItem.Content)[:800])
	}
	err = InsertBestdori(bestdoriChartItem)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

// MysqlSyncToMeiliSearch call at 30s each minute
func MysqlSyncToMeiliSearch() (err error) {
	lastUpdate, err := GetMeiliLastUpdate()
	if err != nil {
		return err
	}
	documents, err := QueryBestdoriFanMadeByLastUpdate(lastUpdate)
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
			err := AddDocument(documents[i:end])
			if err != nil {
				for _, doc := range documents {
					err := AddDocument(doc)
					if err != nil {
						logrus.Log.Warningf("Doc[%d] Add Failed, err: %s", doc.ChartID, err)
						count--
					}
				}
			}
		}
		logrus.Log.Infof("MeiliSearch Sync Finish! %d documents sync!", count)
	} else {
		// Nothing to sync
		logrus.Log.Info("Nothing to do, MeiliSearch Sync Sleep")
	}
	return nil
}

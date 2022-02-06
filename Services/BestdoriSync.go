package Services

import (
	"ayachanV2/Config"
	"ayachanV2/Databases"
	"ayachanV2/Models"
	"ayachanV2/Models/chartFormat"
	"ayachanV2/utils"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type chartDataRequest struct {
	Result bool                          `json:"result"`
	Info   chartFormat.BestdoriChartItem `json:"info"`
}

type chartDataList struct {
	Result bool `json:"result"`
	Count  int  `json:"count"`
	List   []struct {
		ChartID int                `json:"id"`
		Author  chartFormat.Author `json:"author"`
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
		log.Printf("Page %d/%d\n", i, totalPage)
		_, errCode, err := BestdoriFanMadeSyncPage(i)
		if err != nil {
			return errCode, err
		}
	}
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
		log.Printf("Sync Page %d", item)
		_, errCode, err := BestdoriFanMadeSyncPage(item)
		if err != nil {
			return errCode, err
		}
	}
	log.Printf("SyncFinish")
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
			log.Printf("Failed to fetch page info %d [Attempt %d]", page, i)
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
			err = Databases.UpdateBestdori(chartFormat.BestdoriChartUpdateItem{
				ChartID:  item.ChartID,
				Username: item.Author.Username,
				Nickname: item.Author.Nickname,
				Diff:     item.Diff,
				Level:    item.Level,
				Likes:    item.Likes,
			})
		} else {
			go func(i int, item int, ch chan bool) {
				defer func() {
					<-ch
					err := recover()
					if err != nil {
						log.Printf("Panic While Updating Chart #%d : %s", item, err)
					}
					wg.Done()
				}()
				ch <- true
				wg.Add(1)
				var j int
				for j = 1; j <= 5; j++ {
					errCode, err = BestdoriFanMadeInsertID(item)
					if err == nil {
						//log.Printf("Success Update Chart %d [Attemp %d]", item, j)
						break
					} else {
						log.Printf("Failed to update Chart %d [Attemp %d] : Error %s", item, j, err.Error())
					}
				}
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
		log.Printf("谱面无法解析")
		return http.StatusBadRequest, err
	}
	Map := request.Info.Chart.Decode()

	// Insert
	bestdoriChartItem.MapInfoBasic, _, _, _, _ = basicInfoGetter(Map)
	_, bestdoriChartItem.IrregularInfo = ParseMap(Map)

	if bestdoriChartItem.IrregularInfo.Irregular == Models.RegularTypeUnknown {
		log.Printf("分析异常chartID：%d", chartID)
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

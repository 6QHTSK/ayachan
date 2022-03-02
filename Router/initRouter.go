package Router

import (
	"github.com/6QHTSK/ayachan/Config"
	"github.com/6QHTSK/ayachan/Controllers"
	"github.com/6QHTSK/ayachan/Log"
	"github.com/6QHTSK/ayachan/Midware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() (router *gin.Engine) {
	if Config.Config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.New()
	router.Use(Midware.Logger(Log.Log), gin.Recovery())
	router.Use(cors.Default())
	return router
}

func InitAPIV2(router *gin.Engine) {
	v2 := router.Group("/v2")
	{
		v2.GET("/version", Controllers.GetVersion)
		v2.StaticFile("/chart-display", "songList.json")
		// 计算Bestdori谱面信息
		chartInfo := v2.Group("/map-info")
		{
			chartInfo.GET("/bestdori/:chartID", Controllers.MapInfoFromBestdori)
			chartInfo.POST("/", Controllers.MapInfo)
		}
		// 获得爬虫获得的Bestdori信息
		bestdoriInfo := v2.Group("/bestdori")
		{
			bestdoriInfo.GET("/info", Controllers.BestdoriOverAllInfo)
			bestdoriInfo.GET("/list", Controllers.BestdoriFanMadeSearch)
			bestdoriInfo.GET("/list/:chartID", Controllers.BestdoriFanMadeGet)
			bestdoriInfo.GET("/list/:chartID/map-info", Controllers.MapInfoFromBestdori)
			sync := bestdoriInfo.Group("/sync")
			{
				sync.GET("/", Controllers.SyncRand)
				sync.GET("/:chartID", Controllers.SyncChartID)
			}
		}
		// 谱面格式读取和转换
		chartData := v2.Group("/map-data")
		{
			chartData.GET("/bestdori/:chartID", Controllers.MapDataFromBestdori)
			chartData.POST("/", Controllers.MapData)
		}
		sonolus := v2.Group("/sonolus")
		{
			sonolus.POST("/upload/script", Controllers.RedirectSonolusUploadScript)
			sonolus.POST("/upload/song", Controllers.RedirectSonolusUploadSong)
		}
	}
}

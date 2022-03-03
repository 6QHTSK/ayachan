package server

import (
	"github.com/6QHTSK/ayachan/internal/pkg/ginx"
	"github.com/6QHTSK/ayachan/internal/pkg/logrus"
	"github.com/6QHTSK/ayachan/internal/server/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() (router *gin.Engine) {
	if config.Config.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	router = gin.New()
	router.Use(ginx.Logger(logrus.Log), gin.Recovery())
	router.Use(cors.Default())
	return router
}

func InitAPI(router *gin.Engine) {
	v2 := router.Group("/v2")
	{
		v2.GET("/version", GetVersion)
		v2.StaticFile("/chart-display", "songList.json")
		// 计算Bestdori谱面信息
		chartInfo := v2.Group("/map-info")
		{
			chartInfo.GET("/bestdori/:chartID", MapInfoFromBestdori)
			chartInfo.POST("/", MapInfo)
		}
		// 获得爬虫获得的Bestdori信息
		bestdoriInfo := v2.Group("/bestdori")
		{
			bestdoriInfo.GET("/info", BestdoriOverAllInfo)
			bestdoriInfo.GET("/list", BestdoriFanMadeSearch)
			bestdoriInfo.GET("/list/:chartID", BestdoriFanMadeGet)
			bestdoriInfo.GET("/list/:chartID/map-info", MapInfoFromBestdori)
		}
		// 谱面格式读取和转换
		chartData := v2.Group("/map-data")
		{
			chartData.GET("/bestdori/:chartID", MapDataFromBestdori)
			chartData.POST("/", MapData)
		}
		sonolus := v2.Group("/sonolus")
		{
			sonolus.POST("/upload/script", RedirectSonolusUploadScript)
			sonolus.POST("/upload/song", RedirectSonolusUploadSong)
		}
	}
}

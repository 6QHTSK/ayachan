package Router

import (
	"ayachanV2/Controllers"
	_ "ayachanV2/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func InitRouter() (router *gin.Engine) {
	router = gin.Default()
	router.Use(cors.Default())
	return router
}

func InitSwaggerDoc(router *gin.Engine) {
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func InitAPIV2(router *gin.Engine) {
	v2 := router.Group("/v2")
	{
		v2.GET("/version", Controllers.GetVersion)
		// 展示谱面橱窗和详细信息
		//chartDisplay := v2.Group("/chart-display")
		//{
		//	chartDisplay.GET("/")
		//	chartDisplay.GET("/:chartID")
		//}
		v2.StaticFile("/chart-display", "songList.json")
		// 计算Bestdori谱面信息
		chartInfo := v2.Group("/map-info")
		{
			chartInfo.GET("/bestdori/:chartID", Controllers.MapInfoFromBestdori)
			chartInfo.POST("/", Controllers.MapInfo)
		}
		// 获得爬虫获得的Bestdori信息
		bestdoriInfo := v2.Group("/bestdori-info")
		{
			bestdoriInfo.GET("/charter-post-rank", Controllers.CharterPostRank)
			bestdoriInfo.GET("/charter-like-rank", Controllers.CharterLikeRank)
			bestdoriInfo.GET("/song-like-rank", Controllers.SongLikeRank)
			//bestdoriInfo.GET("/charter-list", Controllers.CharterList)
			// TODO BestdoriInfo Add Other API
			/*bestdoriInfoCharter := bestdoriInfo.Group("/charter/:charter")
			{
				bestdoriInfoCharter.GET("/basic-info", Controllers.CharterSelfBasicInfo)
				bestdoriInfoCharter.GET("/post", Controllers.CharterSelfPost)
				bestdoriInfoCharter.GET("/like-rank", Controllers.CharterSelfLikeRank)
				bestdoriInfoCharter.GET("/note-rank", Controllers.CharterSelfNoteRank)
				bestdoriInfoCharter.GET("/time-rank", Controllers.CharterSelfTimeRank)
				bestdoriInfoCharter.GET("/nps-rank", Controllers.CharterSelfNPSRank)
			}*/
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
	}
}

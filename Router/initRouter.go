package Router

import (
	"ayachanV2/Controllers"
	"ayachanV2/Log"
	"ayachanV2/Midware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() (router *gin.Engine) {
	router = gin.New()
	router.Use(Midware.Logger(Log.Log), gin.Recovery())
	router.Use(cors.Default())
	return router
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
		bestdoriInfo := v2.Group("/bestdori")
		{
			bestdoriInfo.GET("/charter-post-rank", Controllers.CharterPostRank)
			bestdoriInfo.GET("/charter-like-rank", Controllers.CharterLikeRank)
			bestdoriInfo.GET("/song-like-rank", Controllers.SongLikeRank)

			bestdoriInfo.GET("/list", Controllers.BestdoriFanMadeSearch)
			bestdoriInfo.GET("/list/:chartID", Controllers.BestdoriFanMadeGet)
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
				//sync.GET("/all", Controllers.SyncAll)
				sync.GET("/:chartID", Controllers.SyncChartID)
				//sync.GET("/refresh-blacklist", Controllers.SyncBlackList)
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

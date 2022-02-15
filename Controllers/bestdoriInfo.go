package Controllers

import (
	"github.com/6QHTSK/ayachan/Databases"
	"github.com/6QHTSK/ayachan/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type InfoOutput struct {
	Result bool        `json:"result"`
	List   interface{} `json:"list,omitempty"`
}

// CharterPostRank 获取谱师发表谱面的排行
//@description 获取谱师发表谱面的排行
//@Summary 获取谱师发表谱面的排行
//@tags BestdoriInfo
//@Param page query int false "页码,默认1"
//@Param limit query int false "页面限制，默认20"
//@Produce json
//@Success 200 {object} InfoOutput "对应的谱面"
//@Failed 500 {object} utils.ErrorObject
//@Router /charter-post-rank [get]
func CharterPostRank(c *gin.Context) {
	page, suc := utils.ConvertQueryInt(c, "page", "1")
	if !suc {
		return
	}
	limit, suc := utils.ConvertQueryInt(c, "limit", "1")
	if !suc {
		return
	}
	list, err := Databases.GetCharterPostRank(page, limit)
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   list,
	})
}

// CharterLikeRank 获取谱师获赞数的排行
//@description 获取谱师获赞数的排行
//@Summary 获取谱师获赞数的排行
//@tags BestdoriInfo
//@Param page query int false "页码,默认1"
//@Param limit query int false "页面限制，默认20"
//@Produce json
//@Success 200 {object} InfoOutput "对应的谱面"
//@Failed 500 {object} utils.ErrorObject
//@Router /charter-like-rank [get]
func CharterLikeRank(c *gin.Context) {
	page, suc := utils.ConvertQueryInt(c, "page", "1")
	if !suc {
		return
	}
	limit, suc := utils.ConvertQueryInt(c, "limit", "1")
	if !suc {
		return
	}
	list, err := Databases.GetCharterLikeRank(page, limit)
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   list,
	})
}

// SongLikeRank 获取谱面获赞数的排行
//@description 获取谱面获赞数的排行
//@Summary 获取谱面获赞数的排行
//@tags BestdoriInfo
//@Param page query int false "页码,默认1"
//@Param limit query int false "页面限制，默认20"
//@Produce json
//@Success 200 {object} InfoOutput "对应的谱面"
//@Failed 500 {object} utils.ErrorObject
//@Router /song-like-rank [get]
func SongLikeRank(c *gin.Context) {
	page, suc := utils.ConvertQueryInt(c, "page", "1")
	if !suc {
		return
	}
	limit, suc := utils.ConvertQueryInt(c, "limit", "1")
	if !suc {
		return
	}
	list, err := Databases.SongLikeRank(page, limit)
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   list,
	})
}

// CharterList 获取所追踪的谱师列表
//@description 获取所追踪的谱师列表
//@Summary 获取所追踪的谱师列表
//@tags BestdoriInfo
//@Produce json
//@Success 200 {object} InfoOutput "谱师列表"
//@Failed 500 {object} utils.ErrorObject
//@Router /charter-list [get]
func CharterList(c *gin.Context) {
	list, err := Databases.GetCharterList()
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   list,
	})
}

// CharterSelfBasicInfo 获取谱师的基本信息
//@description 获取谱师的基本信息
//@Summary 获取谱师的基本信息
//@tags BestdoriInfo
//@Param charter path string true "谱师用户名"
//@Param page query int false "页码,默认1"
//@Param limit query int false "页面限制，默认20"
//@Produce json
//@Success 200 {object} InfoOutput "谱师基础信息"
//@Failed 500 {object} utils.ErrorObject
//@Router /charter/:charter/basic-info [get]
func CharterSelfBasicInfo(c *gin.Context) {
	charter := c.Param("charter")
	info, err := Databases.GetCharterSelfBasic(charter)
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   info,
	})
}

// CharterSelfPost 获取谱师发表的谱面
//@description 获取谱师发表的谱面
//@Summary 获取谱师发表的谱面
//@tags BestdoriInfo
//@Param charter path string true "谱师用户名"
//@Param page query int false "页码,默认1"
//@Param limit query int false "页面限制，默认20"
//@Produce json
//@Success 200 {object} InfoOutput "谱面列表"
//@Failed 500 {object} utils.ErrorObject
//@Router /charter/:charter/post [get]
func CharterSelfPost(c *gin.Context) {
	charter := c.Param("charter")
	page, suc := utils.ConvertQueryInt(c, "page", "1")
	if !suc {
		return
	}
	limit, suc := utils.ConvertQueryInt(c, "limit", "1")
	if !suc {
		return
	}
	list, err := Databases.GetCharterSelfPost(charter, page, limit)
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   list,
	})
}

// CharterSelfLikeRank 获取谱师发表的谱面的赞数排行
//@description 获取谱师发表的赞数排行
//@Summary 获取谱师发表的赞数排行
//@tags BestdoriInfo
//@Param charter path string true "谱师用户名"
//@Param page query int false "页码,默认1"
//@Param limit query int false "页面限制，默认20"
//@Produce json
//@Success 200 {object} InfoOutput "谱面列表"
//@Failed 500 {object} utils.ErrorObject
//@Router /charter/:charter/like-rank [get]
func CharterSelfLikeRank(c *gin.Context) {
	charter := c.Param("charter")
	page, suc := utils.ConvertQueryInt(c, "page", "1")
	if !suc {
		return
	}
	limit, suc := utils.ConvertQueryInt(c, "limit", "1")
	if !suc {
		return
	}
	list, err := Databases.GetCharterSelfLikeRank(charter, page, limit)
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   list,
	})
}

// CharterSelfNoteRank 获取谱师发表的谱面的音符数排行
//@description 获取谱师发表的音符数排行
//@Summary 获取谱师发表的音符数排行
//@tags BestdoriInfo
//@Param charter path string true "谱师用户名"
//@Param page query int false "页码,默认1"
//@Param limit query int false "页面限制，默认20"
//@Produce json
//@Success 200 {object} InfoOutput "谱面列表"
//@Failed 500 {object} utils.ErrorObject
//@Router /charter/:charter/note-rank [get]
func CharterSelfNoteRank(c *gin.Context) {
	charter := c.Param("charter")
	page, suc := utils.ConvertQueryInt(c, "page", "1")
	if !suc {
		return
	}
	limit, suc := utils.ConvertQueryInt(c, "limit", "1")
	if !suc {
		return
	}
	list, err := Databases.GetCharterSelfNoteRank(charter, page, limit)
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   list,
	})
}

// CharterSelfTimeRank 获取谱师发表的谱面的时间排行
//@description 获取谱师发表的谱面的时间排行
//@Summary 获取谱师发表的谱面的时间排行
//@tags BestdoriInfo
//@Param charter path string true "谱师用户名"
//@Param page query int false "页码,默认1"
//@Param limit query int false "页面限制，默认20"
//@Produce json
//@Success 200 {object} InfoOutput "谱面列表"
//@Failed 500 {object} utils.ErrorObject
//@Router /charter/:charter/time-rank [get]
func CharterSelfTimeRank(c *gin.Context) {
	charter := c.Param("charter")
	page, suc := utils.ConvertQueryInt(c, "page", "1")
	if !suc {
		return
	}
	limit, suc := utils.ConvertQueryInt(c, "limit", "1")
	if !suc {
		return
	}
	list, err := Databases.GetCharterSelfTimeRank(charter, page, limit)
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   list,
	})
}

// CharterSelfNPSRank 获取谱师发表的谱面的NPS排行
//@description 获取谱师发表的谱面的NPS排行
//@Summary 获取谱师发表的谱面的NPS排行
//@tags BestdoriInfo
//@Param charter path string true "谱师用户名"
//@Param page query int false "页码,默认1"
//@Param limit query int false "页面限制，默认20"
//@Produce json
//@Success 200 {object} InfoOutput "谱面列表"
//@Failed 500 {object} utils.ErrorObject
//@Router /charter/:charter/nps-rank [get]
func CharterSelfNPSRank(c *gin.Context) {
	charter := c.Param("charter")
	page, suc := utils.ConvertQueryInt(c, "page", "1")
	if !suc {
		return
	}
	limit, suc := utils.ConvertQueryInt(c, "limit", "1")
	if !suc {
		return
	}
	list, err := Databases.GetCharterSelfNPSRank(charter, page, limit)
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   list,
	})
}

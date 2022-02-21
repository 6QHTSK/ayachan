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

func BestdoriOverAllInfo(c *gin.Context) {
	info, err := Databases.GetBestdoriOverallInfo()
	if err != nil {
		utils.ErrorHandle(c, http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, InfoOutput{
		Result: true,
		List:   info,
	})
}

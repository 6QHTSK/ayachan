package Controllers

import (
	"github.com/6QHTSK/ayachan/Models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetVersion 获得该API的版本
//@description 根据内部信息得到API的版本
//@Summary 获得API版本
//@tags Version
//@Produce json
//@Success 200 {object} Models.APIVersion "获得的API版本号"
//@Router /version [get]
func GetVersion(c *gin.Context) {
	var version Models.APIVersion
	version.GetVersion()
	c.JSON(http.StatusOK, version)
}

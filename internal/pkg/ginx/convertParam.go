package ginx

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func ConvertParamInt(c *gin.Context, param string) (number int, success bool) {
	str := c.Param(param)
	number, err := strconv.Atoi(str)
	if err != nil {
		ErrorHandle(c, http.StatusBadRequest, fmt.Errorf("无法解析参数%s", param))
		return number, false
	}
	return number, true
}

func ConvertQueryInt(c *gin.Context, query string, defaultValue string) (number int, success bool) {
	str := c.DefaultQuery(query, defaultValue)
	number, err := strconv.Atoi(str)
	if err != nil {
		ErrorHandle(c, http.StatusBadRequest, fmt.Errorf("无法解析参数%s", query))
		return number, false
	}
	return number, true
}

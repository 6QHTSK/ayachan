package utils

import (
	"ayachanV2/Models/mapFormat"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

type MapInput struct {
	MapFormatIn string            `json:"map_format_in"`
	Options     map[string]string `json:"Options"`
	Map         interface{}       `json:"map"`
}

type BestdoriV2Input struct {
	Map mapFormat.BestdoriV2Chart `json:"map"`
}

type MapOutput struct {
	Result bool        `json:"result"`
	Map    interface{} `json:"map"`
}

// ReadMap 读取内容字段内容 转换为map 会处理错误
func ReadMap(c *gin.Context) (Map mapFormat.Chart, Options map[string]string, suc bool) {
	var inputOptions MapInput
	err := c.ShouldBindBodyWith(&inputOptions, binding.JSON)
	if err != nil {
		ErrorHandle(c, http.StatusBadRequest, fmt.Errorf("无法处理传入谱面"))
		return Map, Options, false
	}

	Options = inputOptions.Options

	if inputOptions.MapFormatIn == "BestdoriV2" {
		var input BestdoriV2Input
		err := c.BindJSON(&input)
		if err != nil {
			ErrorHandle(c, http.StatusBadRequest, fmt.Errorf("无法处理传入谱面"))
			return Map, Options, false
		}
		Map = input.Map.Decode()
	} else {
		ErrorHandle(c, http.StatusBadRequest, fmt.Errorf("不支持的传入格式：%s", inputOptions.MapFormatIn))
		return Map, Options, false
	}

	return Map, Options, true
}

// ReturnMap 返回指定格式的谱面 会处理错误
func ReturnMap(c *gin.Context, Map mapFormat.Chart, formatOut string) (suc bool) {
	if formatOut == "BestdoriV2" {
		BestdoriV2map, err := Map.EncodeBestdoriV2()
		if err != nil {
			ErrorHandle(c, http.StatusInternalServerError, fmt.Errorf("转换内部格式错误：%s", err.Error()))
			return false
		}
		c.JSON(http.StatusOK, MapOutput{
			Result: true,
			Map:    BestdoriV2map,
		})
	} else {
		ErrorHandle(c, http.StatusBadRequest, fmt.Errorf("不支持的传出格式：%s", formatOut))
		return false
	}
	return true
}

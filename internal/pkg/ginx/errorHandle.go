package ginx

import (
	"github.com/gin-gonic/gin"
)

type ErrorObject struct {
	Result  bool   `json:"result"`
	Message string `json:"message"`
}

func ErrorHandle(c *gin.Context, httpCode int, err error) bool {
	if err != nil {
		returnObject := ErrorObject{
			Result:  false,
			Message: err.Error(),
		}
		c.JSON(httpCode, returnObject)
		return true
	}
	return false
}

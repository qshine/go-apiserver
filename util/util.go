package util

import (
	"github.com/gin-gonic/gin"
	"github.com/teris-io/shortid"
)

// 生成一个非有序的短id
func GenShortId() (string, error) {
	return shortid.Generate()
}

func GetReqID(c *gin.Context) string {
	v, ok := c.Get("X-Request-Id")
	if !ok {
		return ""
	}
	if requestId, ok := v.(string); ok {
		return requestId
	}
	return ""
}

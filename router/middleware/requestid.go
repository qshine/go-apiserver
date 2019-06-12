package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/satori/go.uuid"
)

// 在请求和返回的 Header 中插入 X-Request-Id
func RequestId() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 如果请求头中有直接利用, 否则生成
		requestId := c.Request.Header.Get("X-Request-Id")
		if requestId == "" {
			u4, _ := uuid.NewV4()
			requestId = u4.String()
		}

		c.Set("X-Request-Id", requestId)
		// 写入返回包的header中
		c.Writer.Header().Set("X-Request-Id", requestId)
		c.Next()
	}
}

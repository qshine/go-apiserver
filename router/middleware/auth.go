package middleware

import (
	"github.com/gin-gonic/gin"
	"go-apiserver/handler"
	"go-apiserver/pkg/errno"
	"go-apiserver/pkg/token"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := token.ParseRequest(c); err != nil {
			handler.SendResponse(c, errno.ErrTokenInvalid, nil)
			// 如果token验证不通过, 忽略该请求并返回结果
			c.Abort()
			return
		}
		c.Next()
	}
}

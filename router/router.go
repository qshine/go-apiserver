package router

import (
	"github.com/gin-gonic/gin"
	"go-apiserver/handler/sd"
	"go-apiserver/handler/user"
	"go-apiserver/router/middleware"
	"net/http"
)

// 加载中间件, 路由, handlers
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// 使用Use函数来为每一个请求设置Header

	// 在处理某些请求时可能因为程序bug或者其他异常情况导致程序panic
	// 为了不影响下一次请求, 通过设置gin.Recovery()来恢复API服务器
	g.Use(gin.Recovery())
	g.Use(middleware.NoCache)
	g.Use(mw...)

	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// 用户组
	u := g.Group("/v1/user")
	{
		u.POST("", user.Create)
		u.DELETE("/:id", user.Delete)
		u.PUT("/:id", user.Update)
		u.GET("", user.List)
		u.GET("/:username", user.Get)
	}

	// 健康检查handler组. 分别被路由到不同的处理函数
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g

}

package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go-apiserver/demo1/router"
	"log"
	"net/http"
	"time"
)

// 程序入库函数, 主要做配置文件解析,程序初始化, 路由加载

func main() {
	// 创建engine
	g := gin.New()
	middlewares := []gin.HandlerFunc{}

	// 调用router.Load来加载路由
	router.Load(
		g,
		middlewares..., )

	// 启动的时候开一个协程验证是否成功
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Print("The router has been deployed successfully.")
	}()

	log.Printf("Start to listening address: %s...", "8080")
	log.Printf(http.ListenAndServe(":8080", g).Error())
}

// 检查接口正常
func pingServer() error {
	for i := 0; i < 3; i++ {
		resp, err := http.Get("http://127.0.0.1:8080" + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		log.Print("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}

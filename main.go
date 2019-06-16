package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/lexkong/log"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"go-apiserver/config"
	"go-apiserver/model"
	"go-apiserver/router"
	"go-apiserver/router/middleware"
	"net/http"
	"time"
)

var (
	// 命令行指定配置文件, 通过命令行指定不同的配置文件可以方便的区分不同的环境
	cfg = pflag.StringP("config", "c", "", "apiserver config file path")
)

// 程序入库函数, 主要做配置文件解析,程序初始化, 路由加载
func main() {
	pflag.Parse()
	// 初始化配置
	if err := config.Init(*cfg); err != nil {
		panic(err)
	}

	// 连接并最后进行关闭
	model.DB.Init()
	defer model.DB.Close()

	// 设置gin的运行模式 debug/release/test
	gin.SetMode(viper.GetString("runmode"))

	// 创建engine
	g := gin.New()

	// 调用router.Load来加载路由
	router.Load(
		g,
		middleware.Logging(),   // 日志中间件
		middleware.RequestId(), // 请求id中间件
	)

	// 启动的时候开一个协程验证是否成功
	go func() {
		if err := pingServer(); err != nil {
			log.Fatal("The router has no response, or it might took too long to start up.", err)
		}
		log.Info("The router has been deployed successfully.")
	}()

	// 如果提供了http相关的配置则开启https端口
	cert := viper.GetString("tls.cert")
	key := viper.GetString("tls.key")
	if cert != "" && key != "" {
		go func() {
			log.Infof("Start to listening the incoming requests on https address: %s", viper.GetString("tls.addr"))
			log.Info(http.ListenAndServeTLS(viper.GetString("tls.addr"), cert, key, g).Error())
		}()
	}

	log.Infof("Start to listening address: %s...", viper.GetString("addr"))
	log.Info(http.ListenAndServe(viper.GetString("addr"), g).Error())
}

// 检查接口正常
func pingServer() error {
	for i := 0; i < viper.GetInt("max_ping_count"); i++ {
		resp, err := http.Get(viper.GetString("url") + "/sd/health")
		if err == nil && resp.StatusCode == 200 {
			return nil
		}
		log.Info("Waiting for the router, retry in 1 second.")
		time.Sleep(time.Second)
	}
	return errors.New("Cannot connect to the router.")
}

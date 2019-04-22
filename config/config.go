package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/lexkong/log"
	"github.com/spf13/viper"
	"strings"
)

type Config struct {
	Name string
}

// 初始化配置文件
func (c *Config) initConfig() error {
	if c.Name != "" {
		// 如果指定了配置文件则加载
		viper.SetConfigFile(c.Name)
	} else {
		// 否则加载默认配置conf/config.yaml
		viper.AddConfigPath("conf")
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")     // 设置配置文件类型
	viper.AutomaticEnv()            // 读取匹配的环境变量
	viper.SetEnvPrefix("APISERVER") // 设置环境变量名格式为APISERVER开头. 如设置APISERVER_RUNMODE=debug, 用_连接
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// 监控配置文件变化并热更新配置
func (c *Config) watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("Config file changed: %s", e.Name)
	})
}

// 日志初始化需要放在initConfig后面, 因为要利用配置
func (c *Config) initLog() {
	passLagerCfg := log.PassLagerCfg{
		Writers:        viper.GetString("log.writers"),
		LoggerLevel:    viper.GetString("log.logger_level"),
		LoggerFile:     viper.GetString("log.logger_file"),
		LogFormatText:  viper.GetBool("log.log_format_text"),
		RollingPolicy:  viper.GetString("log.rollingPolicy"),
		LogRotateDate:  viper.GetInt("log.log_rotate_date"),
		LogRotateSize:  viper.GetInt("log.log_rotate_size"),
		LogBackupCount: viper.GetInt("log.log_backup_count"),
	}
	log.InitWithConfig(&passLagerCfg)
}

func Init(cfg string) error {
	// 设置配置文件名字
	c := Config{
		Name: cfg,
	}

	// 初始化配置文件
	if err := c.initConfig(); err != nil {
		return err
	}

	// 初始化日志
	c.initLog()

	// 监控配置文件变化并热加载程序
	c.watchConfig()

	return nil
}

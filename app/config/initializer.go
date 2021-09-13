package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

func init() {
	initConfigure()
}

func initConfigure() *viper.Viper {
	v := viper.New()
	v.SetConfigName("coin-watcher-config") // 设置文件名称（无后缀）
	v.SetConfigType("toml")               // 设置后缀名 {"1.6以后的版本可以不设置该后缀"}
	v.AddConfigPath("./config")           // 设置文件所在路径
	v.Set("verbose", true)                // 设置默认参数

	readConfig(v)

	// 监控配置和重新获取配置
	v.OnConfigChange(func(e fsnotify.Event) {
		//fmt.Println("Config file changed:", e.Name)
		readConfig(v)
	})
	v.WatchConfig()
	return v
}

func readConfig(v *viper.Viper) {
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln(" Config file not found; ignore error if desired")
		} else {
			log.Fatalln("Config file was found but another error was produced")
		}
	} else {
		setUpConfig(v)
	}
}

func setUpConfig(v *viper.Viper) {
	BotToken1 = v.GetString("bot.1.token")
	BotHookPath1 = v.GetString("bot.1.hookPath")
	ServeAddr = v.GetString("server.serveAddr")
	TgWebHookHost = v.GetString("server.tgWebHookHost")
	BasePath   = v.GetString("server.basePath")
	TgWebHookPath = v.GetString("server.tgWebHookPath")
	MysqlDsn = v.GetString("mysql.dsn")
	RedisAddr = v.GetString("redis.addr")
	RedisPwd = v.GetString("redis.password")
	EnableSSL = v.GetBool("server.enableSSL")
	CertFile = v.GetString("server.certFile")
	KeyFile = v.GetString("server.keyFile")
}

package main

import (
	"IMCHAT/models"
	"IMCHAT/router"
	"IMCHAT/utils"
	"github.com/spf13/viper"
	"time"
)

func main() {
	//配置文件加载
	utils.InitConfig()
	//数据库初始化
	utils.InitMysql()
	//redis初始化
	utils.InitRedis()
	//InitTimer()

	//路由注册
	r := router.Router()
	//端口注册
	err := r.Run(":8080")
	if err != nil {
		return
	}
}

// InitTimer 定时器
func InitTimer() {
	utils.Timer(time.Duration(viper.GetInt("timeout.DelayHeartbeat"))*time.Second, time.Duration(viper.GetInt("timeout.HeartbeatHz"))*time.Second, models.CleanConnection, "")
}

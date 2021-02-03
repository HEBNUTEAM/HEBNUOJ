package utils

import (
	"github.com/spf13/viper"
	"os"
)

func InitDbConfig() {
	workDir, _ := os.Getwd()
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/config")
	err := viper.ReadInConfig()
	if err != nil {
		Log("config.log", 5).Println("数据库配置文件读取异常")
	}
}

package util

import (
	"CubeUniverse/universalFuncs"
	"github.com/spf13/viper"
)

func InitConfig() {
	workDir := universalFuncs.GetParentDir() + "/control-backend/main" //获取当前目录
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/login-kit/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

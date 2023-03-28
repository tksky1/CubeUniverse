package util

import (
	// "CubeUniverse/universalFuncs"
	"github.com/spf13/viper"
	"os"
)

func InitConfig() {
	// workDir := universalFuncs.GetParentDir() + "/control-backend/main" //获取当前目录
	//测试用，记得删除和复原
	workDir, _ := os.Getwd() //获取当前目录
	viper.SetConfigName("application")
	viper.SetConfigType("yml")
	viper.AddConfigPath(workDir + "/login-kit/config")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

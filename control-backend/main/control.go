package main

import (
	"control-backend/apikit"
	"control-backend/login-kit/common"
	"control-backend/login-kit/model"
	"control-backend/login-kit/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:"varchar(11);not null;unique"`
	Password  string `gorm:"size(255);not null"`
}

func InitUsrAdmin() {
	db := common.GetDB()
	name := "Admin"
	telephone := "12345678901"
	password := "12345678"
	//判断Admin用户是否已经存在
	var user model.User
	db.Where("name=?", name).First(&user)
	//不存在则创建用户
	if user.ID == 0 {
		//创建用户
		hashdPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //密码hash化

		newUser := model.User{
			Name:      name,
			Telephone: telephone,
			Password:  string(hashdPassword),
		}
		if err := db.Create(&newUser).Error; err != nil {
			panic("createUser err" + err.Error())
		}
	}

}

func loginInit() {
	util.InitConfig()
	db := common.GetDB()
	InitUsrAdmin()
	defer db.Close()
	var r *gin.Engine = gin.Default()
	r = CollectRoute(r) //一次性注册完路由

	port := viper.GetString("server.port")
	if port != "" {
		panic(r.Run(":" + port))
	}
	panic(r.Run())

}

func main() {
	//只是测试的时候先执行这个，正常情况下应该先执行apikit的init
	loginInit()
	//实际上应该先执行这个init
	apikit.Init()

	// 后端内容...
	//初始化登录，完成路由注册，实现全部服务
	loginInit()
	//后端完成前先hold
	select {}
}

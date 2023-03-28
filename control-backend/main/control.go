package main

import (
	"CubeUniverse/universalFuncs"
	"control-backend/cubeControl"
	"control-backend/login-kit/common"
	"control-backend/login-kit/model"
	"control-backend/login-kit/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var UUID = uuid.New().String()

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
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //密码hash化

		newUser := model.User{
			Name:      name,
			Telephone: telephone,
			Password:  string(hashedPassword),
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

// 必须先运行这个调用pod互斥锁防止多个pod同时运行功能
func mutexInit() {
	for {
		locked, _, lockTime := universalFuncs.CheckInUse(cubeControl.ClientSet, "backend-mutex")
		if !locked || time.Now().Sub(lockTime).Seconds() > 5 {
			universalFuncs.SetInUse(cubeControl.ClientSet, "backend-mutex", UUID)
			break
		}
		time.Sleep(3 * time.Second)
	}
	// 启动心跳go程
	go universalFuncs.HeartBeat(cubeControl.ClientSet, "backend-mutex", UUID)
}

func main() {

	cubeControl.ClientSet = universalFuncs.GetClientSet()
	cubeControl.DynamicClient = universalFuncs.GetDynamicClient()

	//TODO：删除测试内容
	test()
	//只是测试的时候先执行这个，正常情况下应该先执行cubekit的init
	loginInit()
	//实际上应该先执行这两个init
	mutexInit()
	cubeControl.Init()

	// 后端内容...
	//初始化登录，完成路由注册，实现全部服务
	loginInit()

}

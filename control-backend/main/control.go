package main

import (
	"CubeUniverse/universalFuncs"
	"context"
	"control-backend/cubeControl"
	"control-backend/login-kit/common"
	"control-backend/login-kit/model"
	"control-backend/login-kit/util"
	"fmt"
	"log"
	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 开启的端口号
var portaddress = ":" + viper.GetString("server.port")

type User struct {
	gorm.Model
	Name     string `gorm:"type:varchar(20);not null"`
	Uid      string `gorm:"varchar(11);not null;unique"`
	Password string `gorm:"size(255);not null"`
}

func InitUsrAdmin() {
	db := common.GetDB()
	if db == nil {
		fmt.Println("DB not exist")
		return
	}
	name := "Admin"
	uid := "12345678901"
	password := "cubeuniverse"
	//判断Admin用户是否已经存在
	var user model.User
	db.Where("name=?", name).First(&user)
	//不存在则创建用户
	if user.ID == 0 {
		//创建用户
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) //密码hash化

		newUser := model.User{
			Name:     name,
			Uid:      uid,
			Password: string(hashedPassword),
		}
		if err := db.Create(&newUser).Error; err != nil {
			log.Println("createUser err" + err.Error())
		}
	}
}

func loginInit(ch1 chan bool) {
	util.InitConfig()
	//初始化DB
	common.InitDB()
	db := common.GetDB()

	InitUsrAdmin()
	dbSQL, _ := db.DB()
	defer dbSQL.Close()
	//表明初始化完成
	ch1 <- true
	gin.SetMode(gin.ReleaseMode)

	//保证正好在下一次重开web服务前之前的关闭
	var r *gin.Engine = gin.Default()
	r = CollectRoute(r) //一次性注册完路由

	port := viper.GetString("server.port")

	if port != "" {
		panic(r.Run(":" + port))
	}
	if err := r.Run(); err != nil {
		log.Println("err in init run " + err.Error())
	}

}

// 启用pod互斥锁，必须在初始化时运行
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

var UUID = uuid.New().String()

// 监视数据库初始化状态
func watchDB(ch1 chan bool, srv *http.Server) {
	if ok := <-ch1; ok { //阻塞会一直等待
		//接收到后就关闭旧的监听
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Println("err in wait:   " + err.Error())
		}
		log.Println("wait closed")
		return
	} else {
		log.Println("wait err")
	}
}
func main() {
	gin.SetMode(gin.ReleaseMode)

	cubeControl.ClientSet = universalFuncs.GetClientSet()
	cubeControl.DynamicClient = universalFuncs.GetDynamicClient()
	mutexInit() // 等待pod锁释放

	log.Println("正在加载control-backend..")
	gin.SetMode(gin.ReleaseMode)
	//先开启一个web服务告诉前端需要等待
	router := gin.Default()
	router.GET("/api/wait", func(c *gin.Context) {
		c.String(http.StatusBadRequest, "{\"msg\":\"waiting\"}")
	})

	srv := &http.Server{
		Addr:    portaddress,
		Handler: router,
	}
	go srv.ListenAndServe()
	//建立一个channel来协程通信
	ch1 := make(chan bool)
	//在loginInit正式提供dashbord开始前先检测是否数据库被初始化了
	//开启一个协程来检查
	go watchDB(ch1, srv)

	//测试websocket发送数据

	cubeControl.Init()

	go cubeControl.GetLog()

	// 后端内容...
	//初始化登录，完成路由注册，实现全部服务
	loginInit(ch1)

}

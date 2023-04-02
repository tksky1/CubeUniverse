package common

import (
	"control-backend/login-kit/model"
	"control-backend/login-kit/util"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
)

var DB *gorm.DB

// var DB *gorm.DB //TODO:调试后修改

func InitDB() {

	//TODO:调试加入等待操作，模拟数据库连接很慢10min
	//time.Sleep(30 * time.Second)
	//记得删除

	util.InitConfig()
	driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	defaultDatabase := viper.GetString("defaultdatabase")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		defaultDatabase,
		charset) //集成一个完整的参数用于数据库登录创建新的数据库
	fmt.Println(args)

	db, err := gorm.Open(driverName, args)
	if err != nil {
		fmt.Println(err.Error())
	}
	db.Exec("CREATE DATABASE IF NOT EXISTS " + database)
	args = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		username,
		password,
		host,
		port,
		database,
		charset) //集成一个完整的参数用于数据库登录
	db, err = gorm.Open(driverName, args)

	if err != nil {
		fmt.Println(err.Error())
	}
	db.AutoMigrate(&model.User{}) //自动创建表
	DB = db
}

func GetDB() *gorm.DB {
	return DB
}

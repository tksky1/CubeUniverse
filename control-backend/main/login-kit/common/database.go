package common

import (
	"control-backend/login-kit/model"
	"control-backend/login-kit/util"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/spf13/viper"
)

var DB *gorm.DB

func InitDB() {

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

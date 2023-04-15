package common

import (
	"control-backend/login-kit/model"
	"control-backend/login-kit/util"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() {

	util.InitConfig()
	//driverName := viper.GetString("datasource.driverName")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	//database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	//charset := viper.GetString("datasource.charset")
	defaultDatabase := viper.GetString("datasource.defaultdatabase")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, username, password,
		defaultDatabase, port)
	log.Println(dsn)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN: dsn, PreferSimpleProtocol: true,
	}), &gorm.Config{})

	for err != nil {
		log.Println(err.Error())
		db, err = gorm.Open(postgres.New(postgres.Config{
			DSN: dsn, PreferSimpleProtocol: true,
		}), &gorm.Config{})
	}
	err = db.AutoMigrate(&model.User{}) //自动创建表
	if err != nil {
		log.Println(err.Error())
	}
	DB = db
}

func GetDB() *gorm.DB {
	return DB
}

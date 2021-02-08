package common

import (
	"fmt"
	"github.com/HEBNUOJ/model"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := viper.GetString("datasource.drivername")
	host := viper.GetString("datasource.host")
	port := viper.GetString("datasource.port")
	database := viper.GetString("datasource.database")
	username := viper.GetString("datasource.username")
	password := viper.GetString("datasource.password")
	charset := viper.GetString("datasource.charset")
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		username,
		password,
		host,
		port,
		database,
		charset)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		log.Fatal("failed to connect database ,err:", err)
	}
	db.SingularTable(true) // 禁止表名为结构体的复数
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.LoginLog{})
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}

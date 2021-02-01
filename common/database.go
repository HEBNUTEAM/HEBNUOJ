package common

import (
	"fmt"
	"github.com/HEBNUOJ/model"
	"github.com/jinzhu/gorm"
	"log"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginnessential"
	username := "root"
	password := "69719900"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",
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
	DB = db
	return db
}

func GetDB() *gorm.DB {
	return DB
}

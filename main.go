package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"time"
)

type User struct {
	id         int    `gorm:autoIncrement;primaryKey`
	email      string `gorm:"type:varchar(100);uniqueIndex"`
	submit     int
	solved     int
	defunct    bool      `gorm:"not null"`
	ip         string    `gorm:"type:varchar(48);not null"`
	createTime time.Time `gorm:"type:datetime;autoCreateTime;not null"`
	password   string    `gorm:"type:varchar(20);not null"`
	nickName   string    `gorm:"type:varchar(100);not null"`
	school     string    `gorm:"type:varchar(100);"`
	role       string    `gorm:"type:varchar(20);not null"`
}

func main() {
	db := InitDB()
	defer db.Close()

	r := gin.Default()
	r.GET("/api/auth/register", func(ctx *gin.Context) {

		// 获取参数
		nickname := ctx.PostForm("nick")
		email := ctx.PostForm("email")
		//verification := ctx.PostForm("verification")
		password1 := ctx.PostForm("pwd1")
		password2 := ctx.PostForm("pwd2")

		if len(nickname) > 100 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "用户名的长度必须小于100字节",
			})
		}

		if len(password1) < 6 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "密码不能小于6位",
			})
		}

		if password1 != password2 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "两次密码不一致",
			})
		}
		if isEmailExist(db, email) {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": 422,
				"msg":  "邮箱已存在",
			})
			return
		}

		newUser := User{
			email:      email,
			submit:     0,
			solved:     0,
			defunct:    false,
			ip:         "",
			createTime: time.Now(),
			password:   password1,
			nickName:   nickname,
			school:     "",
			role:       "common",
		}
		db.Create(&newUser)
		ctx.JSON(200, gin.H{
			"msg": "注册成功",
		})
	})
	panic(r.Run())
}

func isEmailExist(db *gorm.DB, email string) bool {
	var user User
	db.Where("email = ?", email).First(&user)
	if len(email) > 0 {
		return true
	}
	return false
}

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
	db.AutoMigrate(&User{})
	return db
}

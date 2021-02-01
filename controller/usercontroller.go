package controller

import (
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/model"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func Register(ctx *gin.Context) {
	db := common.GetDB()
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
		return
	}

	if len(password1) < 6 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "密码不能小于6位",
		})
		return
	}

	if password1 != password2 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "两次密码不一致",
		})
		return
	}
	if isEmailExist(db, email) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "邮箱已存在",
		})
		return
	}

	newUser := model.User{
		Email:      email,
		Submit:     0,
		Solved:     0,
		Defunct:    false,
		Ip:         "",
		CreateTime: time.Now(),
		Password:   password1,
		NickName:   nickname,
		School:     "",
		Role:       "common",
	}
	db.Create(&newUser)
	ctx.JSON(200, gin.H{
		"msg": "注册成功",
	})
}

func isEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	if user.Id > 0 {
		return true
	}
	return false
}

package controller

import (
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/model"
	"github.com/HEBNUOJ/utils"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "加密错误",
		})
		return
	}
	ip := ctx.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	newUser := model.User{
		Email:      email,
		Submit:     0,
		Solved:     0,
		Defunct:    false,
		Ip:         ip,
		CreateTime: time.Now(),
		Password:   string(hasedPassword),
		NickName:   nickname,
		School:     "",
		Role:       "common",
	}
	db.Create(&newUser)
	ctx.JSON(200, gin.H{
		"msg": "注册成功",
	})
}

func Login(ctx *gin.Context) {
	db := common.GetDB()
	// 获取参数
	email := ctx.PostForm("email")
	//verification := ctx.PostForm("verification")
	password := ctx.PostForm("pwd")

	// 判断用户是否存在
	var user model.User
	db.Where("email = ?", email).First(&user)
	if user.Id == 0 {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"code": 422,
			"msg":  "用户不存在或邮箱错误",
		})
		return
	}
	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code": 400,
			"msg":  "密码错误",
		})
	}

	// 发放token给前端
	token, err := common.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code": 500,
			"msg":  "系统异常",
		})
		utils.Log("token.log", 5).Println(err) // 记录错误日志
		return
	}
	ctx.JSON(200, gin.H{
		"code": 200,
		"data": gin.H{"token": token},
		"msg":  "登陆成功",
	})
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	ctx.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{"user": user},
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

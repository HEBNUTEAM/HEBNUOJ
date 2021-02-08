package controller

import (
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/dto"
	"github.com/HEBNUOJ/model"
	"github.com/HEBNUOJ/response"
	"github.com/HEBNUOJ/utils"
	"github.com/HEBNUOJ/vo"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"time"
)

func Register(ctx *gin.Context) {
	db := common.GetDB()
	requestUser := vo.LoginVo{}
	ctx.Bind(&requestUser)
	// 获取参数
	nickname := requestUser.NickName
	email := requestUser.Email
	verification := requestUser.Verification
	captcha := requestUser.Captcha
	captchaId := requestUser.CaptchaId
	password1 := requestUser.Password1
	password2 := requestUser.Password2

	if len(nickname) > 25 || len(nickname) == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil,
			"用户名的长度必须大于等于1个字符，小于等于25个字符")
		return
	}

	if len(password1) < 6 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "密码不能小于6位")

		return
	}

	if password1 != password2 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "两次密码不一致")
		return
	}
	if isEmailExist(db, email) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "邮箱已存在")
		return
	}
	if !isEmailValid(email) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "邮箱不合法")
	}
	if !VerifyCode(captchaId, captcha) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "图形验证码错误")
		return
	}
	if !VerifyEmailCode(email, verification) {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "邮箱验证码错误")
		return
	}
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "加密错误")
		return
	}
	ip := ctx.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}
	// 定义User表字段
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

	// 定义LoginLog表字段
	newLoginLog := model.LoginLog{
		Email:     email,
		Password:  string(hasedPassword),
		Ip:        ip,
		LoginTime: time.Now(),
		Failure:   0,
	}
	db.Create(&newLoginLog)

	response.Success(ctx, nil, "注册成功")
}

func Login(ctx *gin.Context) {
	db := common.GetDB()
	requestUser := vo.LoginVo{}
	ctx.Bind(&requestUser)
	// 获取参数
	email := requestUser.Email
	captchaId := requestUser.CaptchaId
	pngCode := requestUser.Captcha
	password := requestUser.Password1

	var (
		user model.User
		log  model.LoginLog
	)

	// 更新LoginLog表字段
	db.Where("email = ?", email).First(&log)
	hasedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	ip := ctx.ClientIP()
	if ip == "::1" {
		ip = "127.0.0.1"
	}

	log.LoginTime = time.Now()
	log.Password = string(hasedPassword)
	log.Ip = ip

	// 判断用户是否存在
	db.Where("email = ?", email).First(&user)
	if user.Id == 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "用户不存在或邮箱错误")
		return
	}
	if log.Failure > 3 {
		if !VerifyCode(captchaId, pngCode) {
			response.Response(ctx, http.StatusUnprocessableEntity, 422, nil, "图形验证码错误")
			return
		}
	}

	// 判断密码是否正确
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "密码错误")
		log.Failure += 1
		db.Save(&log) // 更新log的全部字段
		return
	}

	// 发放token给前端
	token, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")

		utils.Log("token.log", 5).Println(err) // 记录错误日志
		return
	}

	log.Failure = 0
	db.Save(&log) // 更新log的全部字段

	response.Success(ctx, gin.H{"token": token}, "登陆成功")
}

func Info(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"user": dto.ToUserDto(user.(model.User))}, "")

}

func isEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	if user.Id > 0 {
		return true
	}
	return false
}

func isEmailValid(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// 验证邮箱验证码
func VerifyEmailCode(email, code string) bool {
	client := common.GetRedisClient()
	inCode, err := client.Get(email).Result()
	if err != nil {
		utils.Log("email_code.log", 1).Println("redis get出错", err)
	}
	if inCode == code {
		return true
	}
	return false
}

// 图形验证码验证
func VerifyCode(captchaId, pngCode string) bool {
	if captcha.VerifyString(captchaId, pngCode) {
		return true
	}
	return false
}

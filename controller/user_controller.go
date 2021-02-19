package controller

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/dto"
	"github.com/HEBNUOJ/model"
	"github.com/HEBNUOJ/response"
	"github.com/HEBNUOJ/utils"
	"github.com/HEBNUOJ/vo"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
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

	errString := ""
	switch {
	case len(nickname) > 25 || len(nickname) == 0:
		errString = "用户名的长度必须大于等于1个字符，小于等于25个字符"
	case len(password1) < 6:
		errString = "密码不能小于6位"
	case password1 != password2:
		errString = "两次密码不一致"
	case utils.IsEmailExist(db, email):
		errString = "邮箱已存在"
	case !utils.IsEmailValid(email):
		errString = "邮箱不合法"
	case !utils.VerifyCode(captchaId, captcha):
		errString = "图形验证码错误"
	case !utils.VerifyEmailCode(email, verification):
		errString = "邮箱验证码错误"
	}

	if len(errString) > 0 {
		response.Response(ctx, http.StatusOK, 422, nil, errString)
		return
	}

	// 加密密码
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password1), bcrypt.DefaultCost)
	if err != nil {
		response.Response(ctx, http.StatusOK, 500, nil, "加密错误")
		return
	}

	// 获取ip
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
		Rating:     0,
		Coin:       0,
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
	common.GetRedisClient().Del(ip + ":captcha")
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

	db.Where("email = ?", email).First(&user)
	errCode := 400
	errString := ""

	switch {
	case user.Id == 0: // 判断用户是否存在
		errCode, errString = 422, "用户不存在或邮箱错误"
	case log.Failure > 3 && !utils.VerifyCode(captchaId, pngCode): // 失败次数大于3则需要使用图形验证码
		errCode, errString = 422, "图形验证码错误"
	case !utils.IsPasswordValid(user.Password, password):
		log.Failure += 1
		errString = "密码错误"
		db.Save(&log)
	}

	// 如果出现错误，则返回错误码和错误信息
	if len(errString) > 0 {
		switch {
		case errCode == 400:
			response.Response(ctx, http.StatusOK, 400, nil, "密码错误")
		default:
			response.Response(ctx, http.StatusOK, 422, nil, errString)
		}
		return
	}

	// 将原来的refreshToken删掉， accessToken加入黑名单，ttl设置为10分钟
	jwtToken := ctx.GetHeader("Authorization")
	refreshToken := ctx.GetHeader("RefreshToken")
	common.GetRedisClient().Del(refreshToken)
	common.GetRedisClient().Set(jwtToken, 1, 10*time.Minute)

	// 发放jwtToken给前端
	jwtToken, err := common.ReleaseToken(user)
	if err != nil {
		response.Response(ctx, http.StatusInternalServerError, 500, nil, "系统异常")

		utils.Log("token.log", 5).Println(err) // 记录错误日志
		return
	}
	// 发放refreshToken给前端
	h := md5.New()
	h.Write([]byte(email + strconv.FormatInt(time.Now().Unix(), 10))) // 邮箱和当前时间戳拼接
	cipherStr := h.Sum(nil)
	refreshToken = hex.EncodeToString(cipherStr)

	// 将refreshToken存入redis
	common.GetRedisClient().Set(refreshToken, 1, 72*time.Hour)
	log.Failure = 0
	db.Save(&log)                                // 更新log的全部字段
	common.GetRedisClient().Del(ip + ":captcha") // 清除验证码限制

	response.Success(ctx, gin.H{"jwtToken": jwtToken, "refresh": refreshToken}, "登陆成功")
}

// 退出登录函数
func Logout(ctx *gin.Context) {
	jwtToken := ctx.GetHeader("Authorization")
	refreshToken := ctx.GetHeader("RefreshToken")
	common.GetRedisClient().Del(refreshToken)
	common.GetRedisClient().Set(jwtToken, 1, 10*time.Minute)
	response.Success(ctx, nil, "退出成功")
}

func Info(ctx *gin.Context) {

	user, _ := ctx.Get("user")
	response.Success(ctx, gin.H{"user": dto.ToUserDto(user.(model.User))}, "")
	return
}

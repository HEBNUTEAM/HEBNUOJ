package utils

import (
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/model"
	"github.com/dchest/captcha"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"regexp"
)

func IsEmailExist(db *gorm.DB, email string) bool {
	var user model.User
	db.Where("email = ?", email).First(&user)
	if user.Id > 0 {
		return true
	}
	return false
}

func IsEmailValid(email string) bool {
	pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}

// 验证密码是否有效
func IsPasswordValid(pwd1, pwd2 string) bool {
	if err := bcrypt.CompareHashAndPassword([]byte(pwd1), []byte(pwd2)); err != nil {
		return false
	}
	return true
}

// 验证邮箱验证码
func VerifyEmailCode(email, code string) bool {
	client := common.GetRedisClient()
	inCode, err := client.Get(email).Result()
	if err != nil {
		Log("email_code.log", 1).Println("redis get出错", err)
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

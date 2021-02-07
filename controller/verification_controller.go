package controller

import (
	"bytes"
	"fmt"
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/dto"
	"github.com/HEBNUOJ/response"
	"github.com/HEBNUOJ/utils"
	"github.com/HEBNUOJ/vo"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type CheckCodeController struct{}

// 图形验证码验证
func (serviceCheckCode *CheckCodeController) VerifyCode(ctx *gin.Context) {
	captchaId := ctx.Query("captchaId")
	pngCode := ctx.Query("pngCode")
	if captcha.VerifyString(captchaId, pngCode) {
		response.Success(ctx, nil, "验证成功")
	} else {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "验证码错误")
	}
}

// 加载图形验证码，也可作为初始加载验证码使用，只生成id
func (serviceCheckCode *CheckCodeController) ReloadVerifyCode(ctx *gin.Context) {
	captchaId := captcha.NewLen(4)
	var captcha dto.CaptchaResponse
	captcha.CaptchaId = captchaId
	captcha.ImageUrl = "captcha/" + captchaId + ".png"
	response.Success(ctx, gin.H{"captcha": captcha}, "")
}

// 生成图形验证码
func (serviceCheckCode *CheckCodeController) GenVerifyCode(ctx *gin.Context) {
	// 因为用http传递图片，所以对请求头进行一些设置
	ctx.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Writer.Header().Set("Pragma", "no-cache")
	ctx.Writer.Header().Set("Expires", "0")
	ctx.Writer.Header().Set("Content-Type", "image/png")
	id := ctx.Param("captchaId")
	var content bytes.Buffer
	captcha.WriteImage(&content, id, 100, 50)
	http.ServeContent(ctx.Writer, ctx.Request, id+".png", time.Time{}, bytes.NewReader(content.Bytes()))
}

// 生成邮箱验证码
func (serviceCheckCode *CheckCodeController) GenEmailVerifyCode(ctx *gin.Context) {
	requestUser := vo.LoginVo{}
	ctx.Bind(&requestUser)
	client := common.GetRedisClient()
	email := requestUser.Email
	if client.Exists(email).Val() > 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil,
			"邮箱验证码已存在")
		return
	}
	randCode := make([]byte, 10)
	rand.Read(randCode)
	fmt.Println(randCode)
	err := client.Set(email, randCode, time.Minute*10).Err() // 验证码有效期10分钟
	if err != nil {
		utils.Log("email_code.log", 1).Println("验证码键值对插入redis失败", err)
		return
	}
	sendEmailVerifyCode(email, string(randCode))
	response.Success(ctx, nil, "邮箱验证码申请成功")
}

// 发送邮箱验证码
func sendEmailVerifyCode(email, code string) error {
	mailConn := map[string]string{
		"user": viper.GetString("email.addr"),
		"pass": viper.GetString("email.password"),
		"host": viper.GetString("email.host"),
		"port": viper.GetString("email.port"),
	}
	port, _ := strconv.Atoi(mailConn["port"]) //转换端口类型为int

	m := gomail.NewMessage()
	subject := "测试邮件"
	body := "验证码为：" + code
	m.SetHeader("From", m.FormatAddress(mailConn["user"], "HENUOJ官方"))
	m.SetHeader("To", email)        //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	return err
}

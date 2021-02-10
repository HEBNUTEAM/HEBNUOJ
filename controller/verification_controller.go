package controller

import (
	"bytes"
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/model"
	"github.com/HEBNUOJ/response"
	"github.com/HEBNUOJ/utils"
	"github.com/HEBNUOJ/vo"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
	"net/http"
	"strconv"
	"time"
)

type CheckCodeController struct{}

// 加载图形验证码，也可作为初始加载验证码使用，只生成id
func (serviceCheckCode *CheckCodeController) ReloadVerifyCode(ctx *gin.Context) {
	captchaId := captcha.NewLen(4)
	response.Success(ctx, gin.H{"CaptchaId": captchaId}, "")
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
	email := requestUser.Email
	client := common.GetRedisClient()
	if client.Exists(email).Val() > 0 {
		response.Response(ctx, http.StatusUnprocessableEntity, 422, nil,
			"邮箱验证码已存在")
		return
	}
	code := utils.RandCode(5)
	err := sendEmailVerifyCode(email, code)
	if err != nil { // 如果发送邮件失败则不插入数据库中，防止重复发送
		return
	}
	err = client.Set(email, code, time.Minute*10).Err() // 验证码有效期10分钟
	if err != nil {
		utils.Log("email_code.log", 1).Println("验证码键值对插入redis失败", err)
		return
	}
	response.Success(ctx, nil, "邮箱验证码申请成功")
}

// 判断用户在登陆的时候是否需要验证码
func (serviceCheckCode *CheckCodeController) IsNeedCaptcha(ctx *gin.Context) {
	requestUser := vo.LoginVo{}
	ctx.Bind(&requestUser)
	email := requestUser.Email
	var log = model.LoginLog{}
	common.GetDB().Where("email = ?", email).First(&log)

	if log.Failure > 3 {
		response.Success(ctx, nil, "true")
	} else {
		response.Success(ctx, nil, "false")
	}
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
	body := "验证码为：" + code + ", 本次验证码有效时间为10分钟。"
	m.SetHeader("From", m.FormatAddress(mailConn["user"], "HENUOJ官方"))
	m.SetHeader("To", email)        //发送给多个用户
	m.SetHeader("Subject", subject) //设置邮件主题
	m.SetBody("text/html", body)    //设置邮件正文

	d := gomail.NewDialer(mailConn["host"], port, mailConn["user"], mailConn["pass"])

	err := d.DialAndSend(m)
	if err != nil {
		utils.Log("email_code.log", 1).Println("发送邮件失败", err)
	}
	return err
}

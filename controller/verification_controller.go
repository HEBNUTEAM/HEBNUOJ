package controller

import (
	"bytes"
	"fmt"
	"github.com/HEBNUOJ/dto"
	"github.com/HEBNUOJ/response"
	"github.com/dchest/captcha"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type CheckCodeController struct{}

// 验证码验证
func (serviceCheckCode *CheckCodeController) VerifyCode(ctx *gin.Context) {
	captchaId := ctx.Query("captchaId")
	pngCode := ctx.Query("pngCode")
	if captcha.VerifyString(captchaId, pngCode) {
		response.Success(ctx, nil, "验证成功")
	} else {
		response.Response(ctx, http.StatusBadRequest, 400, nil, "验证码错误")
	}
}

// 加载验证码，也可作为初始加载验证码使用，只生成id
func (serviceCheckCode *CheckCodeController) ReloadVerifyCode(ctx *gin.Context) {
	captchaId := captcha.NewLen(4)
	var captcha dto.CaptchaResponse
	captcha.CaptchaId = captchaId
	captcha.ImageUrl = "captcha/" + captchaId + ".png"
	response.Success(ctx, gin.H{"captcha": captcha}, "")
}

// 生成验证码图片
func (serviceCheckCode *CheckCodeController) GenVerifyCode(ctx *gin.Context) {
	// 因为用http传递图片，所以对请求头进行一些设置
	ctx.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.Writer.Header().Set("Pragma", "no-cache")
	ctx.Writer.Header().Set("Expires", "0")
	ctx.Writer.Header().Set("Content-Type", "image/png")
	id := ctx.Param("id")
	fmt.Println(id)
	var content bytes.Buffer
	captcha.WriteImage(&content, id, 100, 50)
	http.ServeContent(ctx.Writer, ctx.Request, id+".png", time.Time{}, bytes.NewReader(content.Bytes()))
}

package router

import (
	"github.com/HEBNUOJ/controller"
	"github.com/HEBNUOJ/middleware"
	"github.com/gin-gonic/gin"
)

func CollectAuthorizeRoute(r *gin.Engine) *gin.Engine {
	r1 := r.Group("/api/auth") //授权处理路由
	{
		r1.Use(middleware.CorsMiddleware())
		r1.POST("/api/auth/register", controller.Register)
		r1.POST("/api/auth/login", controller.Login)
		r1.POST("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	}
	return r
}

func CollectVerifyRoute(r *gin.Engine) *gin.Engine {
	r1 := r.Group("/api/captcha") // 图形验证码处理路由
	checkCodeController := new(controller.CheckCodeController)
	{
		r1.GET("/refresh", checkCodeController.ReloadVerifyCode)
		r1.GET("/show/:captchaId", checkCodeController.GenVerifyCode)
		r1.POST("/isNeedCaptcha", checkCodeController.IsNeedCaptcha)
	}
	r2 := r.Group("api/email")
	{
		r2.POST("/refresh", checkCodeController.GenEmailVerifyCode)
	}
	return r
}

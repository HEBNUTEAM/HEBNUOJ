package router

import (
	"github.com/HEBNUOJ/controller"
	"github.com/HEBNUOJ/middleware"
	"github.com/gin-gonic/gin"
)

func CollectAuthorizeRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CorsMiddleware())

	r1 := r.Group("/api/auth") //授权处理路由
	{
		r1.POST("/register", controller.Register)
		r1.POST("/login", controller.Login)
		r1.POST("/info", middleware.RenewalTokenMiddleware(),
			middleware.AuthMiddleware(), controller.Info)
		r1.POST("/logout", controller.Logout)
	}
	return r
}

func CollectVerifyRoute(r *gin.Engine) *gin.Engine {
	r1 := r.Group("/api/captcha") // 图形验证码处理路由
	checkCodeController := new(controller.CheckCodeController)
	{
		r1.GET("/refresh", checkCodeController.ReloadVerifyCode)
		r1.POST("/show", checkCodeController.GenVerifyCode)
		r1.POST("/isNeedCaptcha", checkCodeController.IsNeedCaptcha)
	}
	r2 := r.Group("/api/email") // 邮箱验证码处理路由
	{
		r2.POST("/refresh", checkCodeController.GenEmailVerifyCode)
	}
	return r
}

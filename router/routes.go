package router

import (
	"github.com/HEBNUOJ/controller"
	"github.com/HEBNUOJ/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRegisterAndLoginRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CorsMiddleware())
	r.POST("/api/auth/register", controller.Register)
	r.POST("/api/auth/login", controller.Login)
	r.POST("/api/auth/info", middleware.AuthMiddleware(), controller.Info)
	r.POST("/api/auth/")
	return r
}

func CollectVerifyRoute(r *gin.Engine) *gin.Engine {
	//VerifyRoute := r.Group("/verifycode")
	checkCodeController := new(controller.CheckCodeController)
	r.GET("/api/refresh", checkCodeController.ReloadVerifyCode)
	r.GET("/api/show/:id", checkCodeController.GenVerifyCode)
	r.GET("/api/verify", checkCodeController.VerifyCode)
	return r
}

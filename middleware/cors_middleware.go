package middleware

import (
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Add("Access-Control-Allow-Origin", "*")         // 允许的域名
		ctx.Writer.Header().Set("Access-Control-Max-Age", "86400")          // 预检请求的有效期
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "*")        // 设置允许请求的方法
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")        // 设置允许请求的 Header
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // 配置是否可以带认证信息

		//预检请求直接返回200
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(200)
		} else {
			ctx.Next()
		}
	}
}

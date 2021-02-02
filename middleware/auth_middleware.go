package middleware

import (
	"github.com/HEBNUOJ/common"
	"github.com/HEBNUOJ/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 获取authorization header
		tokenString := ctx.GetHeader("Authorization")

		// 验证格式, 判断token是否以"Bearer "为前缀
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort() // 抛弃该次请求
			return
		}

		tokenString = tokenString[7:]
		token, claims, err := common.ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		// 验证通过后获取的claim中的userId
		userID := claims.UserId
		db := common.GetDB()
		var user model.User
		db.First(&user, userID)

		if user.Id == 0 {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "权限不足",
			})
			ctx.Abort()
			return
		}

		// 用户存在，将user信息写入上下文
		ctx.Set("user", user)
		ctx.Next()
	}
}

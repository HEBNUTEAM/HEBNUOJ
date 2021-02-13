package common

import (
	"github.com/HEBNUOJ/model"
	"github.com/dgrijalva/jwt-go"
	"time"
)

var jwtKey = []byte("mdltxdy") // 加盐

type Claims struct {
	UserId int
	jwt.StandardClaims
}

func ReleaseToken(user model.User) (string, error) {
	expirationTime := time.Now().Add(10 * time.Minute)
	claims := &Claims{
		UserId: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(), // token有效期
			IssuedAt:  time.Now().Unix(),     // token发放时间
			Issuer:    "hebnuoj",             // token发放机构
			Subject:   "user token",          // 主题
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 从tokenString中解析出claims
func ParseToken(tokenString string) (*jwt.Token, *Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (i interface{}, err error) {
		return jwtKey, nil
	})
	return token, claims, err
}

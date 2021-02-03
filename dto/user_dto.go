package dto

import (
	"github.com/HEBNUOJ/model"
	"time"
)

type UserDto struct {
	NiceName   string    `json:"nickname"`
	Email      string    `json:"email"`
	Submit     int       `json:"submit"`
	Solved     int       `json:"solved"`
	Defunct    bool      `json:"defunct"`
	CreateTime time.Time `json:"createtime"`
	School     string    `json:"school"`
	Role       string    `json:"role"`
}

func ToUserDto(user model.User) UserDto {
	return UserDto{
		NiceName:   user.NickName,
		Email:      user.Email,
		Submit:     user.Submit,
		Solved:     user.Solved,
		Defunct:    user.Defunct,
		CreateTime: user.CreateTime,
		School:     user.School,
		Role:       user.Role,
	}
}

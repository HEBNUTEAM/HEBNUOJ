package model

import "time"

type User struct {
	Id         int    `gorm:autoIncrement;primaryKey`
	Email      string `gorm:"type:varchar(30);uniqueIndex"`
	Submit     int
	Solved     int
	Defunct    bool      `gorm:"not null"`
	Ip         string    `gorm:"type:varchar(48);not null"`
	CreateTime time.Time `gorm:"type:datetime;autoCreateTime;not null"`
	Password   string    `gorm:"size:60;not null"`
	NickName   string    `gorm:"type:varchar(25);not null"`
	School     string    `gorm:"type:varchar(25);"`
	Role       string    `gorm:"type:varchar(15);not null"`
}

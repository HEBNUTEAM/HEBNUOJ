package model

import "time"

type LoginLog struct {
	Email     string    `gorm:"type:varchar(30);primary_key"`
	Password  string    `gorm:"size:60;not null"`
	Ip        string    `gorm:"type:varchar(48);not null"`
	LoginTime time.Time `gorm:"type:datetime;autoCreateTime;not null"`
	Failure   int       `gorm:"not null"`
}

package model

import "time"

type PublicProblem struct {
	Id           int       `gorm:"autoIncrement;primary_key"`
	Title        string    `gorm:"type:varchar(40);"`
	Description  string    `gorm:"type:text"`
	Input        string    `gorm:"type:text"`
	Output       string    `gorm:"type:text"`
	SampleInput  string    `gorm:"type:text"`
	SampleOutput string    `gorm:"type:text"`
	Spj          bool      `gorm:"default:false;not null"`
	Hint         string    `gorm:"type:text"`
	Source       string    `gorm:"type:varchar(40);"`
	InDate       time.Time `gorm:"type:datetime;"`
	TimeLimit    int       `gorm:"not null"`
	MemoryLimit  int       `gorm:"not null"`
	Defunct      int       `gorm:"default:false;not null"`
	Accepted     int
	Submit       int
	Degree       string `gorm:"type:varchar(10)"`
}

package models

import "gorm.io/gorm"

type User struct {
	gorm.Model           // ID, CreatedAt, UpdatedAt, DeletedAt 자동 포함
	Nickname     string  `gorm:"type:varchar(100);not null"`
	UserID       string  `gorm:"type:varchar(100);not null;unique"`
	Password     string  `gorm:"type:varchar(100);not null"`
	Email        *string `gorm:"type:varchar(100);default:null"` // 포인터를 사용하여 nullable 설정
	VerifiedCode string  `gorm:"type:varchar(100);not null"`
}

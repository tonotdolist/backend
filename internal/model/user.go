package model

import (
	"gorm.io/gorm"
	"tonotdolist/pkg/migrate"
)

func init() {
	migrate.RegisterMigrationModel(&User{})
}

type User struct {
	gorm.Model
	UserId   string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"not null"`
}

func (u *User) TableName() string {
	return "users"
}

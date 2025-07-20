package model

import (
	"gorm.io/gorm"
	"time"
	"tonotdolist/pkg/migrate"
)

func init() {
	migrate.RegisterMigrationModel(&User{})
}

type Activity struct {
	gorm.Model
	UserId      string `gorm:"not null"`
	Type        int8   `gorm:"not null"`
	Name        string `gorm:"not null"`
	Priority    int8
	Description string
	Location    string
	Date        time.Time

	Completed bool `gorm:"default:false"`
}

func (a *Activity) TableName() string {
	return "activities"
}

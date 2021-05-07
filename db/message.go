package db

import (
	"time"

	"gorm.io/gorm"
)

type Message struct {
	gorm.Model
	Phone    string `gorm:"type:varchar(16);not null"`
	Template string `gorm:"type:varchar(8);not null"`
	Salon    uint64 `gorm:"not null"`
	Date     string `gorm:"type:varchar(16);not null"`
	Time     string `gorm:"type:varchar(16);not null"`
	Text     string `gorm:"type:varchar(512)"`
	Done     bool   `gorm:"default:false"`
	Booking  uint64 `gorm:"unique"`
}

type AlarmStatus struct {
	id    uint
	date  *time.Time
	count uint
	salon uint
}

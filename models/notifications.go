package models

import (
	"time"
	"gorm.io/gorm"
)

// Notification represents the notification table
type Notification struct {
	gorm.Model
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	UserID    uint      `gorm:"not null"`
	Message   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	Is_read   bool      `gorm:"default:false"`
}

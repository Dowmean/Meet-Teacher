package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents the user model
type Schedule struct {
	gorm.Model
	ID             uint      `gorm:"primaryKey;autoIncrement"`
	Teacher_id     uint      `json:"Teacherid"`
	Available_time time.Time `json:"Availabletime"`
	Is_booked      bool      `json:"Is_booked"`
	Teacher        User      `gorm:"foreignKey:Teacher_id;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

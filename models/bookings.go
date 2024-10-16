package models

import (
    "time"
	"gorm.io/gorm"
)

type BookingStatus string

const (
    Pending  BookingStatus = "pending"
    Accepted BookingStatus = "accepted"
    Rejected BookingStatus = "rejected"
)

// โครงสร้างข้อมูลสำหรับตาราง Booking
type Booking struct {
	gorm.Model
    ID        uint           `gorm:"primaryKey;autoIncrement"`
    Student_id uint          `gorm:"not null"`
    Teacher_id uint          `gorm:"not null"`
    Time      time.Time      `gorm:"not null"`
    Subject   string         `gorm:"not null"`
    Status    BookingStatus  `gorm:"type:enum('pending', 'accepted', 'rejected');default:'pending'"`
    CreatedAt time.Time
    UpdatedAt time.Time
}
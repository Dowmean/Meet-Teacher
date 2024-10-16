package models

import "gorm.io/gorm"


type Role string

const (
    Teacher Role = "teacher"
    Student Role = "student"
)
// User represents the user model
type User struct {
    gorm.Model
	ID       uint   `gorm:"primaryKey;autoIncrement"`
    Name     string `json:"name"`
    Email    string `json:"email" gorm:"unique"`
    Password string `json:"password"`
	Is_role  Role   `json:"is_role" gorm:"type:enum('teacher','student')"`
	Notifications []Notification `gorm:"foreignKey:UserID"`
}

package models

import "time"

type Attendance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	StudentID uint      `gorm:"index;not null" json:"student_id" binding:"required"`
	Student   User      `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Date      time.Time `gorm:"index;not null" json:"date" binding:"required"`
	Present   bool      `gorm:"default:false" json:"present"`
	MarkedBy  uint      `gorm:"not null" json:"marked_by"`
	Marker    User      `gorm:"foreignKey:MarkedBy" json:"marker,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

//represents attendance statistics
type AttendanceStats struct {
	StudentID            uint    `json:"student_id"`
	PresentDays          int64   `json:"present_days"`
	TotalDays            int64   `json:"total_days"`
	AttendancePercentage float64 `json:"attendance_percentage"`
}

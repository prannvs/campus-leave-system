package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Role string

const (
	RoleAdmin   Role = "admin"
	RoleFaculty Role = "faculty"
	RoleWarden  Role = "warden"
	RoleStudent Role = "student"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"not null" json:"name" binding:"required"`
	Email     string         `gorm:"uniqueIndex;not null" json:"email" binding:"required,email"`
	Password  string         `gorm:"not null" json:"-"`
	Role      Role           `gorm:"type:varchar(20);not null" json:"role" binding:"required"`
	Dept      string         `gorm:"type:varchar(100)" json:"dept"`
	Hostel    string         `gorm:"type:varchar(100)" json:"hostel,omitempty"`
	Leaves    []LeaveRequest `gorm:"foreignKey:StudentID" json:"leaves,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

// creates salted hash
func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

// checks input password with correct password
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

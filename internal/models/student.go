package models

import (
	"time"
)

// Student represents a student in the system
type Student struct {
	ID          uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Email       string       `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Password    string       `gorm:"not null;size:255" json:"-"`
	FullName    string       `gorm:"column:full_name;not null;size:255" json:"full_name"`
	PhoneNumber *string      `gorm:"column:phone_number;size:20" json:"phone_number,omitempty"`
	DateOfBirth *time.Time   `gorm:"column:date_of_birth;type:date" json:"date_of_birth,omitempty"`
	Gender      *string      `gorm:"size:10" json:"gender,omitempty"`
	AvatarURL   *string      `gorm:"column:avatar_url;type:text" json:"avatar_url,omitempty"`
	IsActive    bool         `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt   time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	Enrollments []Enrollment `gorm:"foreignKey:StudentID" json:"-"`
}

// TableName specifies the table name for Student
func (Student) TableName() string {
	return "students"
}

// Valid genders
var ValidGenders = []string{"male", "female", "other"}

// IsValidGender checks if gender is valid
func IsValidGender(gender string) bool {
	for _, g := range ValidGenders {
		if g == gender {
			return true
		}
	}
	return false
}

// IsAdmin checks if student is admin (by email for simplicity)
func (s *Student) IsAdmin() bool {
	return s.Email == "admin@learning.com"
}

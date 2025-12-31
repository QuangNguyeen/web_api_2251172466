package models

import (
	"time"
)

// Lesson represents a lesson in a course
type Lesson struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CourseID    uint      `gorm:"column:course_id;not null" json:"course_id"`
	Title       string    `gorm:"not null;size:255" json:"title"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	VideoURL    *string   `gorm:"column:video_url;type:text" json:"video_url,omitempty"`
	Duration    int       `gorm:"not null;default:0" json:"duration"`
	Order       int       `gorm:"column:order;not null;default:0" json:"order"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Course      *Course   `gorm:"foreignKey:CourseID" json:"-"`
}

// TableName specifies the table name for Lesson
func (Lesson) TableName() string {
	return "lessons"
}

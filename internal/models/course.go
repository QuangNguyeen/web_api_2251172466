package models

import (
	"time"
)

// Course represents a course in the system
type Course struct {
	ID           uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Title        string       `gorm:"not null;size:255" json:"title"`
	Description  *string      `gorm:"type:text" json:"description,omitempty"`
	Instructor   string       `gorm:"not null;size:255" json:"instructor"`
	Category     string       `gorm:"not null;size:50" json:"category"`
	Level        string       `gorm:"not null;size:20" json:"level"`
	Duration     int          `gorm:"not null;default:0" json:"duration"`
	Price        float64      `gorm:"type:decimal(12,2);not null;default:0" json:"price"`
	ImageURL     *string      `gorm:"column:image_url;type:text" json:"image_url,omitempty"`
	Rating       float64      `gorm:"type:decimal(2,1);not null;default:0" json:"rating"`
	StudentCount int          `gorm:"column:student_count;not null;default:0" json:"student_count"`
	LessonCount  int          `gorm:"column:lesson_count;not null;default:0" json:"lesson_count"`
	IsPublished  bool         `gorm:"column:is_published;not null;default:false" json:"is_published"`
	CreatedAt    time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"autoUpdateTime" json:"updated_at"`
	Lessons      []Lesson     `gorm:"foreignKey:CourseID" json:"lessons,omitempty"`
	Enrollments  []Enrollment `gorm:"foreignKey:CourseID" json:"-"`
}

// TableName specifies the table name for Course
func (Course) TableName() string {
	return "courses"
}

// Course categories
var ValidCategories = []string{"Programming", "Design", "Business", "Language", "Music"}

// Course levels
var ValidLevels = []string{"Beginner", "Intermediate", "Advanced"}

// IsValidCategory checks if category is valid
func IsValidCategory(category string) bool {
	for _, c := range ValidCategories {
		if c == category {
			return true
		}
	}
	return false
}

// IsValidLevel checks if level is valid
func IsValidLevel(level string) bool {
	for _, l := range ValidLevels {
		if l == level {
			return true
		}
	}
	return false
}

// Constants for categories and levels
const (
	CategoryProgramming = "Programming"
	CategoryDesign      = "Design"
	CategoryBusiness    = "Business"
	CategoryLanguage    = "Language"
	CategoryMusic       = "Music"

	LevelBeginner     = "Beginner"
	LevelIntermediate = "Intermediate"
	LevelAdvanced     = "Advanced"
)

// Validate validates the course data
func (c *Course) Validate() error {
	if c.Title == "" {
		return ErrValidation
	}
	if c.Instructor == "" {
		return ErrValidation
	}
	if c.Price < 0 {
		return ErrValidation
	}
	if !IsValidCategory(c.Category) {
		return ErrValidation
	}
	if !IsValidLevel(c.Level) {
		return ErrValidation
	}
	return nil
}

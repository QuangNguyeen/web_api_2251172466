package models

import (
	"time"
)

// LessonProgress tracks student progress on individual lessons
type LessonProgress struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	EnrollmentID uint       `gorm:"column:enrollment_id;not null" json:"enrollment_id"`
	LessonID     uint       `gorm:"column:lesson_id;not null" json:"lesson_id"`
	IsCompleted  bool       `gorm:"column:is_completed;not null;default:false" json:"is_completed"`
	CompletedAt  *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
	CreatedAt    time.Time  `gorm:"autoCreateTime" json:"created_at"`

	// Relations
	Enrollment *Enrollment `gorm:"foreignKey:EnrollmentID" json:"-"`
	Lesson     *Lesson     `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
}

// TableName specifies the table name for LessonProgress
func (LessonProgress) TableName() string {
	return "lesson_progress"
}

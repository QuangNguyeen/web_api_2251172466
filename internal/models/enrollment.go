package models

import (
	"time"
)

// Enrollment represents a student's enrollment in a course
type Enrollment struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	StudentID         uint       `gorm:"column:student_id;not null" json:"student_id"`
	CourseID          uint       `gorm:"column:course_id;not null" json:"course_id"`
	EnrollmentDate    time.Time  `gorm:"column:enrollment_date;not null;default:CURRENT_TIMESTAMP" json:"enrollment_date"`
	Progress          int        `gorm:"not null;default:0" json:"progress"`
	Status            string     `gorm:"not null;size:20;default:active" json:"status"`
	LastAccessedAt    *time.Time `gorm:"column:last_accessed_at" json:"last_accessed_at,omitempty"`
	CertificateIssued bool       `gorm:"column:certificate_issued;not null;default:false" json:"certificate_issued"`
	Notes             *string    `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt         time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time  `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Student        *Student         `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Course         *Course          `gorm:"foreignKey:CourseID" json:"course,omitempty"`
	LessonProgress []LessonProgress `gorm:"foreignKey:EnrollmentID" json:"lesson_progress,omitempty"`
}

// TableName specifies the table name for Enrollment
func (Enrollment) TableName() string {
	return "enrollments"
}

// Enrollment statuses
const (
	StatusActive              = "active"
	StatusCompleted           = "completed"
	StatusDropped             = "dropped"
	EnrollmentStatusActive    = "active"
	EnrollmentStatusCompleted = "completed"
	EnrollmentStatusDropped   = "dropped"
)

// ValidEnrollmentStatuses list of valid statuses
var ValidEnrollmentStatuses = []string{StatusActive, StatusCompleted, StatusDropped}

// IsValidEnrollmentStatus checks if status is valid
func IsValidEnrollmentStatus(status string) bool {
	for _, s := range ValidEnrollmentStatuses {
		if s == status {
			return true
		}
	}
	return false
}

// IsActive checks if enrollment is active
func (e *Enrollment) IsActive() bool {
	return e.Status == StatusActive
}

// IsCompleted checks if enrollment is completed
func (e *Enrollment) IsCompleted() bool {
	return e.Status == StatusCompleted
}

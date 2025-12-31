package dto

// EnrollRequest for enrolling in a course
type EnrollRequest struct {
	CourseID uint `json:"course_id" binding:"required"`
}

// EnrollmentResponse for enrollment data
type EnrollmentResponse struct {
	ID                uint             `json:"id"`
	StudentID         uint             `json:"student_id"`
	CourseID          uint             `json:"course_id"`
	EnrollmentDate    string           `json:"enrollment_date"`
	Progress          int              `json:"progress"`
	Status            string           `json:"status"`
	LastAccessedAt    *string          `json:"last_accessed_at,omitempty"`
	CertificateIssued bool             `json:"certificate_issued"`
	Notes             *string          `json:"notes,omitempty"`
	Student           *StudentResponse `json:"student,omitempty"`
	Course            *CourseResponse  `json:"course,omitempty"`
}

// ProgressResponse for detailed progress
type ProgressResponse struct {
	EnrollmentID     uint                     `json:"enrollment_id"`
	Course           *CourseResponse          `json:"course"`
	Progress         int                      `json:"progress"`
	CompletedLessons int                      `json:"completed_lessons"`
	TotalLessons     int                      `json:"total_lessons"`
	Lessons          []LessonProgressResponse `json:"lessons"`
}

// LessonProgressResponse for individual lesson progress
type LessonProgressResponse struct {
	LessonID    uint    `json:"lesson_id"`
	Title       string  `json:"title"`
	IsCompleted bool    `json:"is_completed"`
	CompletedAt *string `json:"completed_at,omitempty"`
}

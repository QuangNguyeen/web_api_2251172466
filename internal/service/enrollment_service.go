package service

import (
	"time"

	"library-api/internal/dto"
	"library-api/internal/models"
	"library-api/internal/repository"

	"gorm.io/gorm"
)

// EnrollmentService handles enrollment business logic
type EnrollmentService struct {
	enrollmentRepo *repository.EnrollmentRepository
	courseRepo     *repository.CourseRepository
	lessonRepo     *repository.LessonRepository
	progressRepo   *repository.LessonProgressRepository
}

// NewEnrollmentService creates a new EnrollmentService
func NewEnrollmentService(
	enrollmentRepo *repository.EnrollmentRepository,
	courseRepo *repository.CourseRepository,
	lessonRepo *repository.LessonRepository,
	progressRepo *repository.LessonProgressRepository,
) *EnrollmentService {
	return &EnrollmentService{
		enrollmentRepo: enrollmentRepo,
		courseRepo:     courseRepo,
		lessonRepo:     lessonRepo,
		progressRepo:   progressRepo,
	}
}

// Enroll enrolls a student in a course
func (s *EnrollmentService) Enroll(studentID uint, req dto.EnrollRequest) (*models.Enrollment, error) {
	// Check if course exists
	course, err := s.courseRepo.FindByIDSimple(req.CourseID)
	if err != nil {
		return nil, models.ErrCourseNotFound
	}

	// Check if already enrolled
	if s.enrollmentRepo.ExistsActiveEnrollment(studentID, req.CourseID) {
		return nil, models.ErrAlreadyEnrolled
	}

	// Get lessons
	lessons, err := s.lessonRepo.FindByCourseID(req.CourseID)
	if err != nil {
		return nil, err
	}

	// Start transaction
	db := s.enrollmentRepo.GetDB()
	tx := db.Begin()

	// Create enrollment
	enrollment := &models.Enrollment{
		StudentID:      studentID,
		CourseID:       req.CourseID,
		EnrollmentDate: time.Now(),
		Progress:       0,
		Status:         models.StatusActive,
	}

	if err := tx.Create(enrollment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Create lesson progress for all lessons
	if len(lessons) > 0 {
		progresses := make([]models.LessonProgress, len(lessons))
		for i, lesson := range lessons {
			progresses[i] = models.LessonProgress{
				EnrollmentID: enrollment.ID,
				LessonID:     lesson.ID,
				IsCompleted:  false,
			}
		}
		if err := tx.Create(&progresses).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// Increment student count
	if err := tx.Model(&models.Course{}).Where("id = ?", req.CourseID).
		UpdateColumn("student_count", course.StudentCount+1).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	enrollment.Course = course
	return enrollment, nil
}

// GetStudentEnrollments gets enrollments for a student
func (s *EnrollmentService) GetStudentEnrollments(studentID uint, status string, page, pageSize int) ([]models.Enrollment, int64, error) {
	return s.enrollmentRepo.FindByStudentID(studentID, status, page, pageSize)
}

// GetCourseEnrollments gets enrollments for a course
func (s *EnrollmentService) GetCourseEnrollments(courseID uint, page, pageSize int) ([]models.Enrollment, int64, error) {
	// Check if course exists
	_, err := s.courseRepo.FindByIDSimple(courseID)
	if err != nil {
		return nil, 0, models.ErrCourseNotFound
	}
	return s.enrollmentRepo.FindByCourseID(courseID, page, pageSize)
}

// CompleteLesson marks a lesson as completed
func (s *EnrollmentService) CompleteLesson(enrollmentID, lessonID, studentID uint, isAdmin bool) (*models.Enrollment, error) {
	// Get enrollment
	enrollment, err := s.enrollmentRepo.FindByIDSimple(enrollmentID)
	if err != nil {
		return nil, models.ErrEnrollmentNotFound
	}

	// Check permission
	if !isAdmin && enrollment.StudentID != studentID {
		return nil, models.ErrForbidden
	}

	// Check status
	if enrollment.Status == models.StatusDropped {
		return nil, models.ErrEnrollmentDropped
	}

	// Get lesson progress
	progress, err := s.progressRepo.FindByEnrollmentAndLesson(enrollmentID, lessonID)
	if err != nil {
		return nil, models.ErrLessonProgressNotFound
	}

	// Check if already completed
	if progress.IsCompleted {
		return nil, models.ErrLessonAlreadyCompleted
	}

	// Start transaction
	db := s.enrollmentRepo.GetDB()
	tx := db.Begin()

	// Mark lesson as completed
	now := time.Now()
	progress.IsCompleted = true
	progress.CompletedAt = &now
	if err := tx.Save(progress).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Calculate new progress
	completed := s.progressRepo.CountCompletedByEnrollmentID(enrollmentID) + 1 // +1 for current
	total := s.progressRepo.CountTotalByEnrollmentID(enrollmentID)
	newProgress := int((float64(completed) / float64(total)) * 100)

	// Update enrollment
	enrollment.Progress = newProgress
	enrollment.LastAccessedAt = &now

	// Check if completed
	if newProgress >= 100 {
		enrollment.Status = models.StatusCompleted
		enrollment.CertificateIssued = true
	}

	if err := tx.Save(enrollment).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return enrollment, nil
}

// GetProgress gets detailed progress for an enrollment
func (s *EnrollmentService) GetProgress(enrollmentID, studentID uint, isAdmin bool) (*dto.ProgressResponse, error) {
	// Get enrollment with course
	enrollment, err := s.enrollmentRepo.FindByID(enrollmentID)
	if err != nil {
		return nil, models.ErrEnrollmentNotFound
	}

	// Check permission
	if !isAdmin && enrollment.StudentID != studentID {
		return nil, models.ErrForbidden
	}

	// Get lesson progress
	progresses, err := s.progressRepo.FindByEnrollmentID(enrollmentID)
	if err != nil {
		return nil, err
	}

	// Build response
	courseResp := ToCourseResponse(enrollment.Course)
	lessonProgresses := make([]dto.LessonProgressResponse, len(progresses))
	completedCount := 0

	for i, p := range progresses {
		lessonProgresses[i] = dto.LessonProgressResponse{
			LessonID:    p.LessonID,
			Title:       p.Lesson.Title,
			IsCompleted: p.IsCompleted,
		}
		if p.CompletedAt != nil {
			at := p.CompletedAt.Format("2006-01-02T15:04:05Z")
			lessonProgresses[i].CompletedAt = &at
		}
		if p.IsCompleted {
			completedCount++
		}
	}

	return &dto.ProgressResponse{
		EnrollmentID:     enrollmentID,
		Course:           &courseResp,
		Progress:         enrollment.Progress,
		CompletedLessons: completedCount,
		TotalLessons:     len(progresses),
		Lessons:          lessonProgresses,
	}, nil
}

// DropEnrollment drops an enrollment
func (s *EnrollmentService) DropEnrollment(enrollmentID, studentID uint, isAdmin bool) error {
	enrollment, err := s.enrollmentRepo.FindByIDSimple(enrollmentID)
	if err != nil {
		return models.ErrEnrollmentNotFound
	}

	// Check permission
	if !isAdmin && enrollment.StudentID != studentID {
		return models.ErrForbidden
	}

	// Check if already dropped
	if enrollment.Status == models.StatusDropped {
		return models.ErrEnrollmentDropped
	}

	// Start transaction
	db := s.enrollmentRepo.GetDB()
	tx := db.Begin()

	// Update status
	enrollment.Status = models.StatusDropped
	if err := tx.Save(enrollment).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Decrement student count
	if err := tx.Model(&models.Course{}).Where("id = ?", enrollment.CourseID).
		UpdateColumn("student_count", gorm.Expr("student_count - 1")).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

// ToEnrollmentResponse converts Enrollment to EnrollmentResponse
func ToEnrollmentResponse(enrollment *models.Enrollment) dto.EnrollmentResponse {
	resp := dto.EnrollmentResponse{
		ID:                enrollment.ID,
		StudentID:         enrollment.StudentID,
		CourseID:          enrollment.CourseID,
		EnrollmentDate:    enrollment.EnrollmentDate.Format("2006-01-02T15:04:05Z"),
		Progress:          enrollment.Progress,
		Status:            enrollment.Status,
		CertificateIssued: enrollment.CertificateIssued,
		Notes:             enrollment.Notes,
	}

	if enrollment.LastAccessedAt != nil {
		at := enrollment.LastAccessedAt.Format("2006-01-02T15:04:05Z")
		resp.LastAccessedAt = &at
	}

	if enrollment.Student != nil {
		studentResp := ToStudentResponse(enrollment.Student)
		resp.Student = &studentResp
	}

	if enrollment.Course != nil {
		courseResp := ToCourseResponse(enrollment.Course)
		resp.Course = &courseResp
	}

	return resp
}

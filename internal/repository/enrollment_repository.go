package repository

import (
	"library-api/internal/models"

	"gorm.io/gorm"
)

// EnrollmentRepository handles enrollment data access
type EnrollmentRepository struct {
	db *gorm.DB
}

// NewEnrollmentRepository creates a new EnrollmentRepository
func NewEnrollmentRepository(db *gorm.DB) *EnrollmentRepository {
	return &EnrollmentRepository{db: db}
}

// Create creates a new enrollment
func (r *EnrollmentRepository) Create(enrollment *models.Enrollment) error {
	return r.db.Create(enrollment).Error
}

// FindByID finds an enrollment by ID with relations
func (r *EnrollmentRepository) FindByID(id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.Preload("Student").Preload("Course").First(&enrollment, id).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// FindByIDSimple finds an enrollment by ID
func (r *EnrollmentRepository) FindByIDSimple(id uint) (*models.Enrollment, error) {
	var enrollment models.Enrollment
	err := r.db.First(&enrollment, id).Error
	if err != nil {
		return nil, err
	}
	return &enrollment, nil
}

// Update updates an enrollment
func (r *EnrollmentRepository) Update(enrollment *models.Enrollment) error {
	return r.db.Save(enrollment).Error
}

// ExistsActiveEnrollment checks if student has active enrollment for course
func (r *EnrollmentRepository) ExistsActiveEnrollment(studentID, courseID uint) bool {
	var count int64
	r.db.Model(&models.Enrollment{}).
		Where("student_id = ? AND course_id = ? AND status != ?", studentID, courseID, models.StatusDropped).
		Count(&count)
	return count > 0
}

// FindByStudentID finds enrollments by student ID with pagination
func (r *EnrollmentRepository) FindByStudentID(studentID uint, status string, page, pageSize int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	query := r.db.Model(&models.Enrollment{}).Where("student_id = ?", studentID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Course").
		Order("enrollment_date DESC").
		Offset(offset).Limit(pageSize).
		Find(&enrollments).Error

	return enrollments, total, err
}

// FindByCourseID finds enrollments by course ID with pagination
func (r *EnrollmentRepository) FindByCourseID(courseID uint, page, pageSize int) ([]models.Enrollment, int64, error) {
	var enrollments []models.Enrollment
	var total int64

	query := r.db.Model(&models.Enrollment{}).Where("course_id = ?", courseID)

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Student").
		Order("enrollment_date DESC").
		Offset(offset).Limit(pageSize).
		Find(&enrollments).Error

	return enrollments, total, err
}

// GetDB returns database instance
func (r *EnrollmentRepository) GetDB() *gorm.DB {
	return r.db
}

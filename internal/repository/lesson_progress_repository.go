package repository

import (
	"library-api/internal/models"

	"gorm.io/gorm"
)

// LessonProgressRepository handles lesson progress data access
type LessonProgressRepository struct {
	db *gorm.DB
}

// NewLessonProgressRepository creates a new LessonProgressRepository
func NewLessonProgressRepository(db *gorm.DB) *LessonProgressRepository {
	return &LessonProgressRepository{db: db}
}

// Create creates a new lesson progress
func (r *LessonProgressRepository) Create(lp *models.LessonProgress) error {
	return r.db.Create(lp).Error
}

// CreateBatch creates multiple lesson progress records
func (r *LessonProgressRepository) CreateBatch(lps []models.LessonProgress) error {
	return r.db.Create(&lps).Error
}

// FindByEnrollmentAndLesson finds progress by enrollment and lesson
func (r *LessonProgressRepository) FindByEnrollmentAndLesson(enrollmentID, lessonID uint) (*models.LessonProgress, error) {
	var lp models.LessonProgress
	err := r.db.Where("enrollment_id = ? AND lesson_id = ?", enrollmentID, lessonID).First(&lp).Error
	if err != nil {
		return nil, err
	}
	return &lp, nil
}

// Update updates lesson progress
func (r *LessonProgressRepository) Update(lp *models.LessonProgress) error {
	return r.db.Save(lp).Error
}

// FindByEnrollmentID finds all progress by enrollment ID
func (r *LessonProgressRepository) FindByEnrollmentID(enrollmentID uint) ([]models.LessonProgress, error) {
	var progresses []models.LessonProgress
	err := r.db.Preload("Lesson").Where("enrollment_id = ?", enrollmentID).Find(&progresses).Error
	return progresses, err
}

// CountCompletedByEnrollmentID counts completed lessons for enrollment
func (r *LessonProgressRepository) CountCompletedByEnrollmentID(enrollmentID uint) int64 {
	var count int64
	r.db.Model(&models.LessonProgress{}).
		Where("enrollment_id = ? AND is_completed = true", enrollmentID).
		Count(&count)
	return count
}

// CountTotalByEnrollmentID counts total lessons for enrollment
func (r *LessonProgressRepository) CountTotalByEnrollmentID(enrollmentID uint) int64 {
	var count int64
	r.db.Model(&models.LessonProgress{}).
		Where("enrollment_id = ?", enrollmentID).
		Count(&count)
	return count
}

// DeleteByEnrollmentID deletes all progress for an enrollment
func (r *LessonProgressRepository) DeleteByEnrollmentID(enrollmentID uint) error {
	return r.db.Where("enrollment_id = ?", enrollmentID).Delete(&models.LessonProgress{}).Error
}

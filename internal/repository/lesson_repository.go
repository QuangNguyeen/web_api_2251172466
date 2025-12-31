package repository

import (
	"library-api/internal/models"

	"gorm.io/gorm"
)

// LessonRepository handles lesson data access
type LessonRepository struct {
	db *gorm.DB
}

// NewLessonRepository creates a new LessonRepository
func NewLessonRepository(db *gorm.DB) *LessonRepository {
	return &LessonRepository{db: db}
}

// Create creates a new lesson
func (r *LessonRepository) Create(lesson *models.Lesson) error {
	return r.db.Create(lesson).Error
}

// FindByID finds a lesson by ID
func (r *LessonRepository) FindByID(id uint) (*models.Lesson, error) {
	var lesson models.Lesson
	err := r.db.First(&lesson, id).Error
	if err != nil {
		return nil, err
	}
	return &lesson, nil
}

// FindByCourseID finds all lessons by course ID
func (r *LessonRepository) FindByCourseID(courseID uint) ([]models.Lesson, error) {
	var lessons []models.Lesson
	err := r.db.Where("course_id = ?", courseID).Order("\"order\" ASC").Find(&lessons).Error
	return lessons, err
}

// CountByCourseID counts lessons by course ID
func (r *LessonRepository) CountByCourseID(courseID uint) int64 {
	var count int64
	r.db.Model(&models.Lesson{}).Where("course_id = ?", courseID).Count(&count)
	return count
}

// Update updates a lesson
func (r *LessonRepository) Update(lesson *models.Lesson) error {
	return r.db.Save(lesson).Error
}

// Delete deletes a lesson
func (r *LessonRepository) Delete(id uint) error {
	return r.db.Delete(&models.Lesson{}, id).Error
}

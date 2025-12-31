package repository

import (
	"library-api/internal/models"

	"gorm.io/gorm"
)

// CourseRepository handles course data access
type CourseRepository struct {
	db *gorm.DB
}

// NewCourseRepository creates a new CourseRepository
func NewCourseRepository(db *gorm.DB) *CourseRepository {
	return &CourseRepository{db: db}
}

// Create creates a new course
func (r *CourseRepository) Create(course *models.Course) error {
	return r.db.Create(course).Error
}

// FindByID finds a course by ID with lessons
func (r *CourseRepository) FindByID(id uint) (*models.Course, error) {
	var course models.Course
	err := r.db.Preload("Lessons", func(db *gorm.DB) *gorm.DB {
		return db.Order("\"order\" ASC")
	}).First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// FindByIDSimple finds a course by ID without lessons
func (r *CourseRepository) FindByIDSimple(id uint) (*models.Course, error) {
	var course models.Course
	err := r.db.First(&course, id).Error
	if err != nil {
		return nil, err
	}
	return &course, nil
}

// Update updates a course
func (r *CourseRepository) Update(course *models.Course) error {
	return r.db.Save(course).Error
}

// Delete deletes a course
func (r *CourseRepository) Delete(id uint) error {
	return r.db.Delete(&models.Course{}, id).Error
}

// FindAll finds all courses with pagination and filters
func (r *CourseRepository) FindAll(page, pageSize int, search, category, level string, minPrice, maxPrice *float64, publishedOnly bool) ([]models.Course, int64, error) {
	var courses []models.Course
	var total int64

	query := r.db.Model(&models.Course{})

	// Search in title, description, instructor
	if search != "" {
		query = query.Where("title ILIKE ? OR description ILIKE ? OR instructor ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Filter by category
	if category != "" {
		query = query.Where("category = ?", category)
	}

	// Filter by level
	if level != "" {
		query = query.Where("level = ?", level)
	}

	// Filter by price range
	if minPrice != nil {
		query = query.Where("price >= ?", *minPrice)
	}
	if maxPrice != nil {
		query = query.Where("price <= ?", *maxPrice)
	}

	// Filter published only
	if publishedOnly {
		query = query.Where("is_published = true")
	}

	// Count total
	query.Count(&total)

	// Pagination
	offset := (page - 1) * pageSize
	query = query.Order("created_at DESC").Offset(offset).Limit(pageSize)

	err := query.Find(&courses).Error
	return courses, total, err
}

// HasActiveEnrollments checks if course has active enrollments
func (r *CourseRepository) HasActiveEnrollments(courseID uint) bool {
	var count int64
	r.db.Model(&models.Enrollment{}).
		Where("course_id = ? AND status != ?", courseID, models.StatusDropped).
		Count(&count)
	return count > 0
}

// IncrementStudentCount increments student_count
func (r *CourseRepository) IncrementStudentCount(courseID uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		UpdateColumn("student_count", gorm.Expr("student_count + 1")).Error
}

// DecrementStudentCount decrements student_count
func (r *CourseRepository) DecrementStudentCount(courseID uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		UpdateColumn("student_count", gorm.Expr("student_count - 1")).Error
}

// IncrementLessonCount increments lesson_count
func (r *CourseRepository) IncrementLessonCount(courseID uint) error {
	return r.db.Model(&models.Course{}).Where("id = ?", courseID).
		UpdateColumn("lesson_count", gorm.Expr("lesson_count + 1")).Error
}

// GetDB returns database instance
func (r *CourseRepository) GetDB() *gorm.DB {
	return r.db
}

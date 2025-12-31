package repository

import (
	"library-api/internal/models"

	"gorm.io/gorm"
)

// StudentRepository handles student data access
type StudentRepository struct {
	db *gorm.DB
}

// NewStudentRepository creates a new StudentRepository
func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// Create creates a new student
func (r *StudentRepository) Create(student *models.Student) error {
	return r.db.Create(student).Error
}

// FindByID finds a student by ID
func (r *StudentRepository) FindByID(id uint) (*models.Student, error) {
	var student models.Student
	err := r.db.First(&student, id).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// FindByEmail finds a student by email
func (r *StudentRepository) FindByEmail(email string) (*models.Student, error) {
	var student models.Student
	err := r.db.Where("email = ?", email).First(&student).Error
	if err != nil {
		return nil, err
	}
	return &student, nil
}

// ExistsByEmail checks if email exists
func (r *StudentRepository) ExistsByEmail(email string) bool {
	var count int64
	r.db.Model(&models.Student{}).Where("email = ?", email).Count(&count)
	return count > 0
}

// Update updates a student
func (r *StudentRepository) Update(student *models.Student) error {
	return r.db.Save(student).Error
}

// FindAll finds all students with pagination and search
func (r *StudentRepository) FindAll(page, pageSize int, search string) ([]models.Student, int64, error) {
	var students []models.Student
	var total int64

	query := r.db.Model(&models.Student{})

	// Search
	if search != "" {
		query = query.Where("email ILIKE ? OR full_name ILIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// Count total
	query.Count(&total)

	// Pagination
	offset := (page - 1) * pageSize
	query = query.Order("created_at DESC").Offset(offset).Limit(pageSize)

	err := query.Find(&students).Error
	return students, total, err
}

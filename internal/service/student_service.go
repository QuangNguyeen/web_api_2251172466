package service

import (
	"time"

	"library-api/internal/dto"
	"library-api/internal/models"
	"library-api/internal/repository"
)

// StudentService handles student business logic
type StudentService struct {
	studentRepo *repository.StudentRepository
}

// NewStudentService creates a new StudentService
func NewStudentService(studentRepo *repository.StudentRepository) *StudentService {
	return &StudentService{
		studentRepo: studentRepo,
	}
}

// GetAll gets all students with pagination
func (s *StudentService) GetAll(page, pageSize int, search string) ([]models.Student, int64, error) {
	return s.studentRepo.FindAll(page, pageSize, search)
}

// GetByID gets a student by ID
func (s *StudentService) GetByID(id uint) (*models.Student, error) {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return nil, models.ErrStudentNotFound
	}
	return student, nil
}

// Update updates a student
func (s *StudentService) Update(id uint, req dto.UpdateStudentRequest, isAdmin bool) (*models.Student, error) {
	student, err := s.studentRepo.FindByID(id)
	if err != nil {
		return nil, models.ErrStudentNotFound
	}

	// Update fields
	if req.FullName != nil {
		student.FullName = *req.FullName
	}
	if req.PhoneNumber != nil {
		student.PhoneNumber = req.PhoneNumber
	}
	if req.DateOfBirth != nil {
		if *req.DateOfBirth != "" {
			parsed, err := time.Parse("2006-01-02", *req.DateOfBirth)
			if err == nil {
				student.DateOfBirth = &parsed
			}
		} else {
			student.DateOfBirth = nil
		}
	}
	if req.Gender != nil {
		if *req.Gender != "" && !models.IsValidGender(*req.Gender) {
			return nil, models.ErrInvalidGender
		}
		student.Gender = req.Gender
	}
	if req.AvatarURL != nil {
		student.AvatarURL = req.AvatarURL
	}
	// Only admin can update is_active
	if req.IsActive != nil && isAdmin {
		student.IsActive = *req.IsActive
	}

	if err := s.studentRepo.Update(student); err != nil {
		return nil, err
	}

	return student, nil
}

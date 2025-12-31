package service

import (
	"time"

	"library-api/internal/config"
	"library-api/internal/dto"
	"library-api/internal/models"
	"library-api/internal/repository"
	"library-api/internal/utils"
)

// MSSV - Mã Sinh Viên
const StudentIDMSSV = "2251172466"

// AuthService handles authentication logic
type AuthService struct {
	studentRepo *repository.StudentRepository
	cfg         *config.Config
}

// NewAuthService creates a new AuthService
func NewAuthService(studentRepo *repository.StudentRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		studentRepo: studentRepo,
		cfg:         cfg,
	}
}

// Register registers a new student
func (s *AuthService) Register(req dto.RegisterRequest) (*models.Student, error) {
	// Check if email already exists
	if s.studentRepo.ExistsByEmail(req.Email) {
		return nil, models.ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Parse date of birth
	var dob *time.Time
	if req.DateOfBirth != nil && *req.DateOfBirth != "" {
		parsed, err := time.Parse("2006-01-02", *req.DateOfBirth)
		if err == nil {
			dob = &parsed
		}
	}

	// Create student
	student := &models.Student{
		Email:       req.Email,
		Password:    hashedPassword,
		FullName:    req.FullName,
		PhoneNumber: req.PhoneNumber,
		DateOfBirth: dob,
		Gender:      req.Gender,
		IsActive:    true,
	}

	if err := s.studentRepo.Create(student); err != nil {
		return nil, err
	}

	return student, nil
}

// Login authenticates a student and returns a token
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find student by email
	student, err := s.studentRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, models.ErrInvalidCredentials
	}

	// Check if active
	if !student.IsActive {
		return nil, models.ErrStudentInactive
	}

	// Check password
	if !utils.CheckPassword(req.Password, student.Password) {
		return nil, models.ErrInvalidCredentials
	}

	// Generate token
	token, err := utils.GenerateToken(student.ID, student.Email, s.cfg.JWTSecret, s.cfg.JWTExpiration, student.IsAdmin())
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:     token,
		StudentID: StudentIDMSSV, // MSSV hardcoded
		User:      ToStudentResponse(student),
	}, nil
}

// GetStudentByID gets a student by ID
func (s *AuthService) GetStudentByID(id uint) (*models.Student, error) {
	return s.studentRepo.FindByID(id)
}

// ToStudentResponse converts Student to StudentResponse
func ToStudentResponse(student *models.Student) dto.StudentResponse {
	resp := dto.StudentResponse{
		ID:          student.ID,
		Email:       student.Email,
		FullName:    student.FullName,
		PhoneNumber: student.PhoneNumber,
		Gender:      student.Gender,
		AvatarURL:   student.AvatarURL,
		IsActive:    student.IsActive,
		CreatedAt:   student.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if student.DateOfBirth != nil {
		dob := student.DateOfBirth.Format("2006-01-02")
		resp.DateOfBirth = &dob
	}

	return resp
}

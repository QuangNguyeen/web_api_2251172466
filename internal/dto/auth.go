package dto

// RegisterRequest for student registration
type RegisterRequest struct {
	Email       string  `json:"email" binding:"required,email"`
	Password    string  `json:"password" binding:"required,min=6"`
	FullName    string  `json:"full_name" binding:"required"`
	PhoneNumber *string `json:"phone_number"`
	DateOfBirth *string `json:"date_of_birth"` // Format: "2006-01-02"
	Gender      *string `json:"gender" binding:"omitempty,oneof=male female other"`
}

// LoginRequest for student login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse for login response
type LoginResponse struct {
	Token     string          `json:"token"`
	StudentID string          `json:"student_id"` // MSSV - hardcoded
	User      StudentResponse `json:"user"`
}

// StudentResponse for student data
type StudentResponse struct {
	ID          uint    `json:"id"`
	Email       string  `json:"email"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number,omitempty"`
	DateOfBirth *string `json:"date_of_birth,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	AvatarURL   *string `json:"avatar_url,omitempty"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
}

// UpdateStudentRequest for updating student
type UpdateStudentRequest struct {
	FullName    *string `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
	DateOfBirth *string `json:"date_of_birth"`
	Gender      *string `json:"gender" binding:"omitempty,oneof=male female other"`
	AvatarURL   *string `json:"avatar_url"`
	IsActive    *bool   `json:"is_active"` // Admin only
}

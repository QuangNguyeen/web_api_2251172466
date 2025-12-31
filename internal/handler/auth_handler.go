package handler

import (
	"net/http"

	"library-api/internal/dto"
	"library-api/internal/models"
	"library-api/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler handles authentication endpoints
type AuthHandler struct {
	authService *service.AuthService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// Register godoc
// @Summary Register new student
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Register request"
// @Success 201 {object} dto.APIResponse{data=dto.StudentResponse}
// @Failure 400 {object} dto.APIResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	student, err := h.authService.Register(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrEmailAlreadyExists {
			status = http.StatusConflict
		}
		c.JSON(status, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Registration successful",
		Data:    service.ToStudentResponse(student),
	})
}

// Login godoc
// @Summary Login student
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login request"
// @Success 200 {object} dto.APIResponse{data=dto.LoginResponse}
// @Failure 401 {object} dto.APIResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		status := http.StatusUnauthorized
		if err == models.ErrStudentInactive {
			status = http.StatusForbidden
		}
		c.JSON(status, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// GetMe godoc
// @Summary Get current student info
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} dto.APIResponse{data=dto.StudentResponse}
// @Failure 401 {object} dto.APIResponse
// @Router /api/auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	studentID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.APIResponse{
			Success: false,
			Message: "Unauthorized",
		})
		return
	}

	student, err := h.authService.GetStudentByID(studentID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: "Student not found",
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    service.ToStudentResponse(student),
	})
}

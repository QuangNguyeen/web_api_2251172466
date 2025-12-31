package handler

import (
	"net/http"
	"strconv"

	"library-api/internal/dto"
	"library-api/internal/models"
	"library-api/internal/service"

	"github.com/gin-gonic/gin"
)

// StudentHandler handles student endpoints
type StudentHandler struct {
	studentService *service.StudentService
}

// NewStudentHandler creates a new StudentHandler
func NewStudentHandler(studentService *service.StudentService) *StudentHandler {
	return &StudentHandler{studentService: studentService}
}

// GetAll godoc
// @Summary Get all students (admin only)
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param search query string false "Search by email or name"
// @Success 200 {object} dto.APIResponse{data=dto.ListResponse[dto.StudentResponse]}
// @Router /api/students [get]
func (h *StudentHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	students, total, err := h.studentService.GetAll(page, limit, search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	items := make([]dto.StudentResponse, len(students))
	for i, s := range students {
		items[i] = service.ToStudentResponse(&s)
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data: dto.ListResponse[dto.StudentResponse]{
			Items:       items,
			TotalItems:  total,
			CurrentPage: page,
			PageSize:    limit,
			TotalPages:  (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetByID godoc
// @Summary Get student by ID
// @Tags Students
// @Security BearerAuth
// @Produce json
// @Param id path int true "Student ID"
// @Success 200 {object} dto.APIResponse{data=dto.StudentResponse}
// @Router /api/students/{id} [get]
func (h *StudentHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid student ID",
		})
		return
	}

	// Check permission: admin or self
	currentUserID, _ := c.Get("userID")
	isAdmin, _ := c.Get("isAdmin")

	if !isAdmin.(bool) && currentUserID.(uint) != uint(id) {
		c.JSON(http.StatusForbidden, dto.APIResponse{
			Success: false,
			Message: "Access denied",
		})
		return
	}

	student, err := h.studentService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    service.ToStudentResponse(student),
	})
}

// Update godoc
// @Summary Update student
// @Tags Students
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Student ID"
// @Param request body dto.UpdateStudentRequest true "Update request"
// @Success 200 {object} dto.APIResponse{data=dto.StudentResponse}
// @Router /api/students/{id} [put]
func (h *StudentHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid student ID",
		})
		return
	}

	// Check permission: admin or self
	currentUserID, _ := c.Get("userID")
	isAdmin, _ := c.Get("isAdmin")

	if !isAdmin.(bool) && currentUserID.(uint) != uint(id) {
		c.JSON(http.StatusForbidden, dto.APIResponse{
			Success: false,
			Message: "Access denied",
		})
		return
	}

	var req dto.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	student, err := h.studentService.Update(uint(id), req, isAdmin.(bool))
	if err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrStudentNotFound {
			status = http.StatusNotFound
		} else if err == models.ErrInvalidGender {
			status = http.StatusBadRequest
		}
		c.JSON(status, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Student updated successfully",
		Data:    service.ToStudentResponse(student),
	})
}

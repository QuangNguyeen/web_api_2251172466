package handler

import (
	"net/http"
	"strconv"

	"library-api/internal/dto"
	"library-api/internal/models"
	"library-api/internal/service"

	"github.com/gin-gonic/gin"
)

// EnrollmentHandler handles enrollment endpoints
type EnrollmentHandler struct {
	enrollmentService *service.EnrollmentService
}

// NewEnrollmentHandler creates a new EnrollmentHandler
func NewEnrollmentHandler(enrollmentService *service.EnrollmentService) *EnrollmentHandler {
	return &EnrollmentHandler{enrollmentService: enrollmentService}
}

// Enroll godoc
// @Summary Enroll in a course
// @Tags Enrollments
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.EnrollRequest true "Enroll request"
// @Success 201 {object} dto.APIResponse{data=dto.EnrollmentResponse}
// @Router /api/enrollments [post]
func (h *EnrollmentHandler) Enroll(c *gin.Context) {
	studentID, _ := c.Get("userID")

	var req dto.EnrollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	enrollment, err := h.enrollmentService.Enroll(studentID.(uint), req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrCourseNotFound {
			status = http.StatusNotFound
		} else if err == models.ErrAlreadyEnrolled {
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
		Message: "Enrolled successfully",
		Data:    service.ToEnrollmentResponse(enrollment),
	})
}

// GetStudentEnrollments godoc
// @Summary Get student enrollments
// @Tags Enrollments
// @Security BearerAuth
// @Produce json
// @Param id path int true "Student ID"
// @Param status query string false "Filter by status"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} dto.APIResponse{data=dto.ListResponse[dto.EnrollmentResponse]}
// @Router /api/students/{id}/enrollments [get]
func (h *EnrollmentHandler) GetStudentEnrollments(c *gin.Context) {
	studentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid student ID",
		})
		return
	}

	// Check permission
	currentUserID, _ := c.Get("userID")
	isAdmin, _ := c.Get("isAdmin")

	if !isAdmin.(bool) && currentUserID.(uint) != uint(studentID) {
		c.JSON(http.StatusForbidden, dto.APIResponse{
			Success: false,
			Message: "Access denied",
		})
		return
	}

	status := c.Query("status")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	enrollments, total, err := h.enrollmentService.GetStudentEnrollments(uint(studentID), status, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	items := make([]dto.EnrollmentResponse, len(enrollments))
	for i, e := range enrollments {
		items[i] = service.ToEnrollmentResponse(&e)
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data: dto.ListResponse[dto.EnrollmentResponse]{
			Items:       items,
			TotalItems:  total,
			CurrentPage: page,
			PageSize:    limit,
			TotalPages:  (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetCourseEnrollments godoc
// @Summary Get course enrollments (admin only)
// @Tags Enrollments
// @Security BearerAuth
// @Produce json
// @Param id path int true "Course ID"
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Success 200 {object} dto.APIResponse{data=dto.ListResponse[dto.EnrollmentResponse]}
// @Router /api/courses/{id}/enrollments [get]
func (h *EnrollmentHandler) GetCourseEnrollments(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid course ID",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	enrollments, total, err := h.enrollmentService.GetCourseEnrollments(uint(courseID), page, limit)
	if err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrCourseNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	items := make([]dto.EnrollmentResponse, len(enrollments))
	for i, e := range enrollments {
		items[i] = service.ToEnrollmentResponse(&e)
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data: dto.ListResponse[dto.EnrollmentResponse]{
			Items:       items,
			TotalItems:  total,
			CurrentPage: page,
			PageSize:    limit,
			TotalPages:  (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// CompleteLesson godoc
// @Summary Mark lesson as completed
// @Tags Enrollments
// @Security BearerAuth
// @Success 200 {object} dto.APIResponse{data=dto.EnrollmentResponse}
// @Router /api/enrollments/{id}/lessons/{lesson_id}/complete [put]
func (h *EnrollmentHandler) CompleteLesson(c *gin.Context) {
	enrollmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid enrollment ID",
		})
		return
	}

	lessonID, err := strconv.ParseUint(c.Param("lesson_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid lesson ID",
		})
		return
	}

	studentID, _ := c.Get("userID")
	isAdmin, _ := c.Get("isAdmin")

	enrollment, err := h.enrollmentService.CompleteLesson(uint(enrollmentID), uint(lessonID), studentID.(uint), isAdmin.(bool))
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case models.ErrEnrollmentNotFound, models.ErrLessonProgressNotFound:
			status = http.StatusNotFound
		case models.ErrForbidden:
			status = http.StatusForbidden
		case models.ErrEnrollmentDropped, models.ErrLessonAlreadyCompleted:
			status = http.StatusConflict
		}
		c.JSON(status, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Lesson completed",
		Data:    service.ToEnrollmentResponse(enrollment),
	})
}

// GetProgress godoc
// @Summary Get enrollment progress
// @Tags Enrollments
// @Security BearerAuth
// @Produce json
// @Param id path int true "Enrollment ID"
// @Success 200 {object} dto.APIResponse{data=dto.ProgressResponse}
// @Router /api/enrollments/{id}/progress [get]
func (h *EnrollmentHandler) GetProgress(c *gin.Context) {
	enrollmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid enrollment ID",
		})
		return
	}

	studentID, _ := c.Get("userID")
	isAdmin, _ := c.Get("isAdmin")

	progress, err := h.enrollmentService.GetProgress(uint(enrollmentID), studentID.(uint), isAdmin.(bool))
	if err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrEnrollmentNotFound {
			status = http.StatusNotFound
		} else if err == models.ErrForbidden {
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
		Data:    progress,
	})
}

// DropEnrollment godoc
// @Summary Drop enrollment
// @Tags Enrollments
// @Security BearerAuth
// @Param id path int true "Enrollment ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/enrollments/{id} [delete]
func (h *EnrollmentHandler) DropEnrollment(c *gin.Context) {
	enrollmentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid enrollment ID",
		})
		return
	}

	studentID, _ := c.Get("userID")
	isAdmin, _ := c.Get("isAdmin")

	err = h.enrollmentService.DropEnrollment(uint(enrollmentID), studentID.(uint), isAdmin.(bool))
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case models.ErrEnrollmentNotFound:
			status = http.StatusNotFound
		case models.ErrForbidden:
			status = http.StatusForbidden
		case models.ErrEnrollmentDropped:
			status = http.StatusConflict
		}
		c.JSON(status, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Enrollment dropped successfully",
	})
}

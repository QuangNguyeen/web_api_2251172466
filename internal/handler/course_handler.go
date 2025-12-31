package handler

import (
	"net/http"
	"strconv"

	"library-api/internal/dto"
	"library-api/internal/models"
	"library-api/internal/service"

	"github.com/gin-gonic/gin"
)

// CourseHandler handles course endpoints
type CourseHandler struct {
	courseService *service.CourseService
}

// NewCourseHandler creates a new CourseHandler
func NewCourseHandler(courseService *service.CourseService) *CourseHandler {
	return &CourseHandler{courseService: courseService}
}

// GetAll godoc
// @Summary Get all courses
// @Tags Courses
// @Produce json
// @Param page query int false "Page" default(1)
// @Param limit query int false "Limit" default(10)
// @Param search query string false "Search"
// @Param category query string false "Category"
// @Param level query string false "Level"
// @Param min_price query number false "Min price"
// @Param max_price query number false "Max price"
// @Param published_only query bool false "Published only" default(true)
// @Success 200 {object} dto.APIResponse{data=dto.ListResponse[dto.CourseResponse]}
// @Router /api/courses [get]
func (h *CourseHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	category := c.Query("category")
	level := c.Query("level")
	publishedOnly := c.DefaultQuery("published_only", "true") == "true"

	var minPrice, maxPrice *float64
	if min := c.Query("min_price"); min != "" {
		if v, err := strconv.ParseFloat(min, 64); err == nil {
			minPrice = &v
		}
	}
	if max := c.Query("max_price"); max != "" {
		if v, err := strconv.ParseFloat(max, 64); err == nil {
			maxPrice = &v
		}
	}

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	courses, total, err := h.courseService.GetAll(page, limit, search, category, level, minPrice, maxPrice, publishedOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	items := make([]dto.CourseResponse, len(courses))
	for i, course := range courses {
		items[i] = service.ToCourseResponse(&course)
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data: dto.ListResponse[dto.CourseResponse]{
			Items:       items,
			TotalItems:  total,
			CurrentPage: page,
			PageSize:    limit,
			TotalPages:  (total + int64(limit) - 1) / int64(limit),
		},
	})
}

// GetByID godoc
// @Summary Get course by ID with lessons
// @Tags Courses
// @Produce json
// @Param id path int true "Course ID"
// @Success 200 {object} dto.APIResponse{data=dto.CourseResponse}
// @Router /api/courses/{id} [get]
func (h *CourseHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid course ID",
		})
		return
	}

	course, err := h.courseService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Data:    service.ToCourseResponse(course),
	})
}

// Create godoc
// @Summary Create course (admin only)
// @Tags Courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body dto.CreateCourseRequest true "Create request"
// @Success 201 {object} dto.APIResponse{data=dto.CourseResponse}
// @Router /api/courses [post]
func (h *CourseHandler) Create(c *gin.Context) {
	var req dto.CreateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	course, err := h.courseService.Create(req)
	if err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrInvalidCategory || err == models.ErrInvalidLevel {
			status = http.StatusBadRequest
		}
		c.JSON(status, dto.APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Course created successfully",
		Data:    service.ToCourseResponse(course),
	})
}

// Update godoc
// @Summary Update course (admin only)
// @Tags Courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param request body dto.UpdateCourseRequest true "Update request"
// @Success 200 {object} dto.APIResponse{data=dto.CourseResponse}
// @Router /api/courses/{id} [put]
func (h *CourseHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid course ID",
		})
		return
	}

	var req dto.UpdateCourseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	course, err := h.courseService.Update(uint(id), req)
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

	c.JSON(http.StatusOK, dto.APIResponse{
		Success: true,
		Message: "Course updated successfully",
		Data:    service.ToCourseResponse(course),
	})
}

// Delete godoc
// @Summary Delete course (admin only)
// @Tags Courses
// @Security BearerAuth
// @Param id path int true "Course ID"
// @Success 200 {object} dto.APIResponse
// @Router /api/courses/{id} [delete]
func (h *CourseHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid course ID",
		})
		return
	}

	err = h.courseService.Delete(uint(id))
	if err != nil {
		status := http.StatusInternalServerError
		if err == models.ErrCourseNotFound {
			status = http.StatusNotFound
		} else if err == models.ErrCourseHasEnrollments {
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
		Message: "Course deleted successfully",
	})
}

// AddLesson godoc
// @Summary Add lesson to course (admin only)
// @Tags Courses
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Course ID"
// @Param request body dto.CreateLessonRequest true "Lesson request"
// @Success 201 {object} dto.APIResponse{data=dto.LessonResponse}
// @Router /api/courses/{id}/lessons [post]
func (h *CourseHandler) AddLesson(c *gin.Context) {
	courseID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid course ID",
		})
		return
	}

	var req dto.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.APIResponse{
			Success: false,
			Message: "Invalid request: " + err.Error(),
		})
		return
	}

	lesson, err := h.courseService.AddLesson(uint(courseID), req)
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

	c.JSON(http.StatusCreated, dto.APIResponse{
		Success: true,
		Message: "Lesson added successfully",
		Data:    service.ToLessonResponse(lesson),
	})
}

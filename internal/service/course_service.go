package service

import (
	"library-api/internal/dto"
	"library-api/internal/models"
	"library-api/internal/repository"
)

// CourseService handles course business logic
type CourseService struct {
	courseRepo *repository.CourseRepository
	lessonRepo *repository.LessonRepository
}

// NewCourseService creates a new CourseService
func NewCourseService(courseRepo *repository.CourseRepository, lessonRepo *repository.LessonRepository) *CourseService {
	return &CourseService{
		courseRepo: courseRepo,
		lessonRepo: lessonRepo,
	}
}

// GetAll gets all courses with filters
func (s *CourseService) GetAll(page, pageSize int, search, category, level string, minPrice, maxPrice *float64, publishedOnly bool) ([]models.Course, int64, error) {
	return s.courseRepo.FindAll(page, pageSize, search, category, level, minPrice, maxPrice, publishedOnly)
}

// GetByID gets a course by ID with lessons
func (s *CourseService) GetByID(id uint) (*models.Course, error) {
	course, err := s.courseRepo.FindByID(id)
	if err != nil {
		return nil, models.ErrCourseNotFound
	}
	return course, nil
}

// Create creates a new course
func (s *CourseService) Create(req dto.CreateCourseRequest) (*models.Course, error) {
	if !models.IsValidCategory(req.Category) {
		return nil, models.ErrInvalidCategory
	}
	if !models.IsValidLevel(req.Level) {
		return nil, models.ErrInvalidLevel
	}

	course := &models.Course{
		Title:       req.Title,
		Description: req.Description,
		Instructor:  req.Instructor,
		Category:    req.Category,
		Level:       req.Level,
		Duration:    req.Duration,
		Price:       req.Price,
		ImageURL:    req.ImageURL,
		IsPublished: false,
	}

	if err := s.courseRepo.Create(course); err != nil {
		return nil, err
	}

	return course, nil
}

// Update updates a course
func (s *CourseService) Update(id uint, req dto.UpdateCourseRequest) (*models.Course, error) {
	course, err := s.courseRepo.FindByIDSimple(id)
	if err != nil {
		return nil, models.ErrCourseNotFound
	}

	// Update fields
	if req.Title != nil {
		course.Title = *req.Title
	}
	if req.Description != nil {
		course.Description = req.Description
	}
	if req.Instructor != nil {
		course.Instructor = *req.Instructor
	}
	if req.Category != nil {
		if !models.IsValidCategory(*req.Category) {
			return nil, models.ErrInvalidCategory
		}
		course.Category = *req.Category
	}
	if req.Level != nil {
		if !models.IsValidLevel(*req.Level) {
			return nil, models.ErrInvalidLevel
		}
		course.Level = *req.Level
	}
	if req.Duration != nil {
		course.Duration = *req.Duration
	}
	if req.Price != nil {
		course.Price = *req.Price
	}
	if req.ImageURL != nil {
		course.ImageURL = req.ImageURL
	}
	if req.IsPublished != nil {
		course.IsPublished = *req.IsPublished
	}

	if err := s.courseRepo.Update(course); err != nil {
		return nil, err
	}

	return course, nil
}

// Delete deletes a course
func (s *CourseService) Delete(id uint) error {
	_, err := s.courseRepo.FindByIDSimple(id)
	if err != nil {
		return models.ErrCourseNotFound
	}

	// Check for active enrollments
	if s.courseRepo.HasActiveEnrollments(id) {
		return models.ErrCourseHasEnrollments
	}

	return s.courseRepo.Delete(id)
}

// AddLesson adds a lesson to a course
func (s *CourseService) AddLesson(courseID uint, req dto.CreateLessonRequest) (*models.Lesson, error) {
	_, err := s.courseRepo.FindByIDSimple(courseID)
	if err != nil {
		return nil, models.ErrCourseNotFound
	}

	lesson := &models.Lesson{
		CourseID:    courseID,
		Title:       req.Title,
		Description: req.Description,
		VideoURL:    req.VideoURL,
		Duration:    req.Duration,
		Order:       req.Order,
	}

	if err := s.lessonRepo.Create(lesson); err != nil {
		return nil, err
	}

	// Update lesson count
	s.courseRepo.IncrementLessonCount(courseID)

	return lesson, nil
}

// ToCourseResponse converts Course to CourseResponse
func ToCourseResponse(course *models.Course) dto.CourseResponse {
	resp := dto.CourseResponse{
		ID:           course.ID,
		Title:        course.Title,
		Description:  course.Description,
		Instructor:   course.Instructor,
		Category:     course.Category,
		Level:        course.Level,
		Duration:     course.Duration,
		Price:        course.Price,
		ImageURL:     course.ImageURL,
		Rating:       course.Rating,
		StudentCount: course.StudentCount,
		LessonCount:  course.LessonCount,
		IsPublished:  course.IsPublished,
		CreatedAt:    course.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if len(course.Lessons) > 0 {
		resp.Lessons = make([]dto.LessonResponse, len(course.Lessons))
		for i, lesson := range course.Lessons {
			resp.Lessons[i] = ToLessonResponse(&lesson)
		}
	}

	return resp
}

// ToLessonResponse converts Lesson to LessonResponse
func ToLessonResponse(lesson *models.Lesson) dto.LessonResponse {
	return dto.LessonResponse{
		ID:          lesson.ID,
		CourseID:    lesson.CourseID,
		Title:       lesson.Title,
		Description: lesson.Description,
		VideoURL:    lesson.VideoURL,
		Duration:    lesson.Duration,
		Order:       lesson.Order,
		CreatedAt:   lesson.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}
}

package dto

// CreateCourseRequest for creating a course
type CreateCourseRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	Instructor  string  `json:"instructor" binding:"required"`
	Category    string  `json:"category" binding:"required,oneof=Programming Design Business Language Music"`
	Level       string  `json:"level" binding:"required,oneof=Beginner Intermediate Advanced"`
	Duration    int     `json:"duration" binding:"min=0"`
	Price       float64 `json:"price" binding:"min=0"`
	ImageURL    *string `json:"image_url"`
}

// UpdateCourseRequest for updating a course
type UpdateCourseRequest struct {
	Title       *string  `json:"title"`
	Description *string  `json:"description"`
	Instructor  *string  `json:"instructor"`
	Category    *string  `json:"category" binding:"omitempty,oneof=Programming Design Business Language Music"`
	Level       *string  `json:"level" binding:"omitempty,oneof=Beginner Intermediate Advanced"`
	Duration    *int     `json:"duration" binding:"omitempty,min=0"`
	Price       *float64 `json:"price" binding:"omitempty,min=0"`
	ImageURL    *string  `json:"image_url"`
	IsPublished *bool    `json:"is_published"`
}

// CourseResponse for course data
type CourseResponse struct {
	ID           uint             `json:"id"`
	Title        string           `json:"title"`
	Description  *string          `json:"description,omitempty"`
	Instructor   string           `json:"instructor"`
	Category     string           `json:"category"`
	Level        string           `json:"level"`
	Duration     int              `json:"duration"`
	Price        float64          `json:"price"`
	ImageURL     *string          `json:"image_url,omitempty"`
	Rating       float64          `json:"rating"`
	StudentCount int              `json:"student_count"`
	LessonCount  int              `json:"lesson_count"`
	IsPublished  bool             `json:"is_published"`
	CreatedAt    string           `json:"created_at"`
	Lessons      []LessonResponse `json:"lessons,omitempty"`
}

// CreateLessonRequest for creating a lesson
type CreateLessonRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	VideoURL    *string `json:"video_url"`
	Duration    int     `json:"duration" binding:"min=0"`
	Order       int     `json:"order" binding:"min=0"`
}

// LessonResponse for lesson data
type LessonResponse struct {
	ID          uint    `json:"id"`
	CourseID    uint    `json:"course_id"`
	Title       string  `json:"title"`
	Description *string `json:"description,omitempty"`
	VideoURL    *string `json:"video_url,omitempty"`
	Duration    int     `json:"duration"`
	Order       int     `json:"order"`
	CreatedAt   string  `json:"created_at"`
}

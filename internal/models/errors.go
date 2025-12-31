package models

import "errors"

// Common errors
var (
	// Student errors
	ErrStudentNotFound    = errors.New("student not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidGender      = errors.New("invalid gender")
	ErrStudentInactive    = errors.New("student account is inactive")

	// Course errors
	ErrCourseNotFound       = errors.New("course not found")
	ErrInvalidCategory      = errors.New("invalid category")
	ErrInvalidLevel         = errors.New("invalid level")
	ErrCourseHasEnrollments = errors.New("cannot delete course with active enrollments")

	// Lesson errors
	ErrLessonNotFound = errors.New("lesson not found")

	// Enrollment errors
	ErrEnrollmentNotFound  = errors.New("enrollment not found")
	ErrAlreadyEnrolled     = errors.New("student already enrolled in this course")
	ErrEnrollmentDropped   = errors.New("enrollment has been dropped")
	ErrEnrollmentCompleted = errors.New("enrollment already completed")
	ErrInvalidStatus       = errors.New("invalid enrollment status")

	// Progress errors
	ErrLessonProgressNotFound = errors.New("lesson progress not found")
	ErrLessonAlreadyCompleted = errors.New("lesson already completed")

	// Auth errors
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	// Validation errors
	ErrValidation = errors.New("validation error")
)

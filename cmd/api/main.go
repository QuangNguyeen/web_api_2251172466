package main

import (
	"log"

	"library-api/internal/config"
	"library-api/internal/handler"
	"library-api/internal/middleware"
	"library-api/internal/repository"
	"library-api/internal/service"

	_ "library-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Online Learning Management API
// @version 1.0
// @description API for Online Learning Management System
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db := config.InitDB(cfg)

	// Initialize repositories
	studentRepo := repository.NewStudentRepository(db)
	courseRepo := repository.NewCourseRepository(db)
	lessonRepo := repository.NewLessonRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	progressRepo := repository.NewLessonProgressRepository(db)

	// Initialize services
	authService := service.NewAuthService(studentRepo, cfg)
	studentService := service.NewStudentService(studentRepo)
	courseService := service.NewCourseService(courseRepo, lessonRepo)
	enrollmentService := service.NewEnrollmentService(enrollmentRepo, courseRepo, lessonRepo, progressRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	studentHandler := handler.NewStudentHandler(studentService)
	courseHandler := handler.NewCourseHandler(courseService)
	enrollmentHandler := handler.NewEnrollmentHandler(enrollmentService)

	// Setup router
	router := gin.Default()

	// Apply CORS middleware
	router.Use(middleware.CORS())

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(cfg.JWTSecret), authHandler.GetMe)
		}

		// Student routes
		students := api.Group("/students")
		students.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			students.GET("", middleware.AdminMiddleware(), studentHandler.GetAll) // Admin only
			students.GET("/:id", studentHandler.GetByID)
			students.PUT("/:id", studentHandler.Update)
			students.GET("/:id/enrollments", enrollmentHandler.GetStudentEnrollments)
		}

		// Course routes
		courses := api.Group("/courses")
		{
			// Public
			courses.GET("", courseHandler.GetAll)
			courses.GET("/:id", courseHandler.GetByID)

			// Admin only
			courses.POST("", middleware.AuthMiddleware(cfg.JWTSecret), middleware.AdminMiddleware(), courseHandler.Create)
			courses.PUT("/:id", middleware.AuthMiddleware(cfg.JWTSecret), middleware.AdminMiddleware(), courseHandler.Update)
			courses.DELETE("/:id", middleware.AuthMiddleware(cfg.JWTSecret), middleware.AdminMiddleware(), courseHandler.Delete)
			courses.POST("/:id/lessons", middleware.AuthMiddleware(cfg.JWTSecret), middleware.AdminMiddleware(), courseHandler.AddLesson)
			courses.GET("/:id/enrollments", middleware.AuthMiddleware(cfg.JWTSecret), middleware.AdminMiddleware(), enrollmentHandler.GetCourseEnrollments)
		}

		// Enrollment routes
		enrollments := api.Group("/enrollments")
		enrollments.Use(middleware.AuthMiddleware(cfg.JWTSecret))
		{
			enrollments.POST("", enrollmentHandler.Enroll)
			enrollments.GET("/:id/progress", enrollmentHandler.GetProgress)
			enrollments.PUT("/:id/lessons/:lesson_id/complete", enrollmentHandler.CompleteLesson)
			enrollments.DELETE("/:id", enrollmentHandler.DropEnrollment)
		}
	}

	// Start server
	log.Println("")
	log.Println("========================================")
	log.Println("Online Learning Management API")
	log.Println("Server: http://localhost:" + cfg.Port)
	log.Println("Swagger: http://localhost:" + cfg.Port + "/swagger/index.html")
	log.Println("========================================")

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

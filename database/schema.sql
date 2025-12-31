-- ============================================
-- Online Learning Management System - Database Schema
-- PostgreSQL 15+
-- Database: db_exam_2251172466
-- ============================================

-- Drop existing tables if they exist (for clean reinstallation)
DROP TABLE IF EXISTS lesson_progress CASCADE;
DROP TABLE IF EXISTS lessons CASCADE;
DROP TABLE IF EXISTS enrollments CASCADE;
DROP TABLE IF EXISTS courses CASCADE;
DROP TABLE IF EXISTS students CASCADE;

-- ============================================
-- Table: students
-- ============================================

CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(20),
    date_of_birth DATE,
    gender VARCHAR(10) CHECK (gender IN ('male', 'female', 'other')),
    avatar_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_students_email ON students(email);

-- ============================================
-- Table: courses
-- ============================================

CREATE TABLE courses (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    instructor VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('Programming', 'Design', 'Business', 'Language', 'Music')),
    level VARCHAR(20) NOT NULL CHECK (level IN ('Beginner', 'Intermediate', 'Advanced')),
    duration INTEGER NOT NULL DEFAULT 0,
    price DECIMAL(12, 2) NOT NULL DEFAULT 0,
    image_url TEXT,
    rating DECIMAL(2, 1) NOT NULL DEFAULT 0.0,
    student_count INTEGER NOT NULL DEFAULT 0,
    lesson_count INTEGER NOT NULL DEFAULT 0,
    is_published BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_courses_category ON courses(category);
CREATE INDEX idx_courses_level ON courses(level);
CREATE INDEX idx_courses_published ON courses(is_published);

-- ============================================
-- Table: lessons
-- ============================================

CREATE TABLE lessons (
    id BIGSERIAL PRIMARY KEY,
    course_id BIGINT NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    video_url TEXT,
    duration INTEGER NOT NULL DEFAULT 0,
    "order" INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_lessons_course 
        FOREIGN KEY (course_id) 
        REFERENCES courses(id) 
        ON DELETE CASCADE
);

CREATE INDEX idx_lessons_course_id ON lessons(course_id);

-- ============================================
-- Table: enrollments
-- ============================================

CREATE TABLE enrollments (
    id BIGSERIAL PRIMARY KEY,
    student_id BIGINT NOT NULL,
    course_id BIGINT NOT NULL,
    enrollment_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    progress INTEGER NOT NULL DEFAULT 0 CHECK (progress >= 0 AND progress <= 100),
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'completed', 'dropped')),
    last_accessed_at TIMESTAMP,
    certificate_issued BOOLEAN NOT NULL DEFAULT false,
    notes TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_enrollments_student 
        FOREIGN KEY (student_id) 
        REFERENCES students(id) 
        ON DELETE RESTRICT,
    
    CONSTRAINT fk_enrollments_course 
        FOREIGN KEY (course_id) 
        REFERENCES courses(id) 
        ON DELETE RESTRICT
);

CREATE INDEX idx_enrollments_student_id ON enrollments(student_id);
CREATE INDEX idx_enrollments_course_id ON enrollments(course_id);
CREATE INDEX idx_enrollments_status ON enrollments(status);

-- Unique constraint: student can only have one active/completed enrollment per course
CREATE UNIQUE INDEX idx_enrollments_unique_active 
    ON enrollments(student_id, course_id) 
    WHERE status != 'dropped';

-- ============================================
-- Table: lesson_progress
-- ============================================

CREATE TABLE lesson_progress (
    id BIGSERIAL PRIMARY KEY,
    enrollment_id BIGINT NOT NULL,
    lesson_id BIGINT NOT NULL,
    is_completed BOOLEAN NOT NULL DEFAULT false,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_lesson_progress_enrollment 
        FOREIGN KEY (enrollment_id) 
        REFERENCES enrollments(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT fk_lesson_progress_lesson 
        FOREIGN KEY (lesson_id) 
        REFERENCES lessons(id) 
        ON DELETE CASCADE,
    
    CONSTRAINT uq_lesson_progress UNIQUE (enrollment_id, lesson_id)
);

CREATE INDEX idx_lesson_progress_enrollment_id ON lesson_progress(enrollment_id);

-- ============================================
-- Trigger: Update updated_at on modification
-- ============================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_students_updated_at 
    BEFORE UPDATE ON students
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_courses_updated_at 
    BEFORE UPDATE ON courses
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_lessons_updated_at 
    BEFORE UPDATE ON lessons
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_enrollments_updated_at 
    BEFORE UPDATE ON enrollments
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

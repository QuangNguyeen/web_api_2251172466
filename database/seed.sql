-- ============================================
-- Online Learning Management System - Seed Data
-- PostgreSQL 15+
-- ============================================

-- Clear existing data
TRUNCATE TABLE lesson_progress CASCADE;
TRUNCATE TABLE lessons CASCADE;
TRUNCATE TABLE enrollments CASCADE;
TRUNCATE TABLE courses CASCADE;
TRUNCATE TABLE students CASCADE;

-- Reset sequences
ALTER SEQUENCE students_id_seq RESTART WITH 1;
ALTER SEQUENCE courses_id_seq RESTART WITH 1;
ALTER SEQUENCE lessons_id_seq RESTART WITH 1;
ALTER SEQUENCE enrollments_id_seq RESTART WITH 1;
ALTER SEQUENCE lesson_progress_id_seq RESTART WITH 1;

-- ============================================
-- Seed Students (5 students)
-- Password: "password123" (BCrypt hashed)
-- ============================================

INSERT INTO students (email, password, full_name, phone_number, date_of_birth, gender, is_active) VALUES
('admin@learning.com', '$2a$10$SSeheAIuC1gpnRwfWiK48OLq1PsPentTmli3FuXynr/gsUYpaOz2a', 'Admin User', '0901234567', '1990-01-15', 'male', true),
('student1@learning.com', '$2a$10$SSeheAIuC1gpnRwfWiK48OLq1PsPentTmli3FuXynr/gsUYpaOz2a', 'Nguyễn Văn An', '0902345678', '2000-05-20', 'male', true),
('student2@learning.com', '$2a$10$SSeheAIuC1gpnRwfWiK48OLq1PsPentTmli3FuXynr/gsUYpaOz2a', 'Trần Thị Bình', '0903456789', '2001-08-10', 'female', true),
('student3@learning.com', '$2a$10$SSeheAIuC1gpnRwfWiK48OLq1PsPentTmli3FuXynr/gsUYpaOz2a', 'Lê Văn Cường', '0904567890', '1999-12-05', 'male', true),
('student4@learning.com', '$2a$10$SSeheAIuC1gpnRwfWiK48OLq1PsPentTmli3FuXynr/gsUYpaOz2a', 'Phạm Thị Dung', '0905678901', '2002-03-25', 'female', true);

-- ============================================
-- Seed Courses (12 courses)
-- ============================================

INSERT INTO courses (title, description, instructor, category, level, duration, price, rating, student_count, lesson_count, is_published) VALUES
-- Programming (4)
('Flutter Development Masterclass', 'Learn Flutter from scratch to build beautiful mobile apps', 'Nguyễn Văn A', 'Programming', 'Beginner', 40, 1500000, 4.8, 150, 20, true),
('Advanced Python Programming', 'Master Python with advanced concepts and real projects', 'Trần Văn B', 'Programming', 'Advanced', 50, 2000000, 4.9, 200, 25, true),
('JavaScript Full Course', 'Complete JavaScript guide for web development', 'Lê Văn C', 'Programming', 'Intermediate', 35, 1200000, 4.7, 180, 18, true),
('React Native Development', 'Build cross-platform mobile apps with React Native', 'Phạm Văn D', 'Programming', 'Intermediate', 45, 1800000, 4.6, 120, 22, true),

-- Design (3)
('UI/UX Design Fundamentals', 'Learn the basics of user interface and experience design', 'Nguyễn Thị E', 'Design', 'Beginner', 30, 1000000, 4.5, 100, 15, true),
('Adobe Photoshop Masterclass', 'Master Photoshop for graphic design', 'Trần Thị F', 'Design', 'Intermediate', 25, 800000, 4.4, 80, 12, true),
('Figma for Beginners', 'Design beautiful interfaces with Figma', 'Lê Thị G', 'Design', 'Beginner', 20, 600000, 4.6, 90, 10, true),

-- Business (2)
('Digital Marketing Complete Guide', 'Learn digital marketing strategies', 'Phạm Văn H', 'Business', 'Beginner', 35, 1500000, 4.7, 130, 16, true),
('Business Analytics with Excel', 'Master Excel for business analytics', 'Nguyễn Văn I', 'Business', 'Intermediate', 28, 900000, 4.3, 70, 14, true),

-- Language (2)
('English for IELTS', 'Prepare for IELTS exam with comprehensive lessons', 'Trần Thị K', 'Language', 'Intermediate', 60, 2500000, 4.9, 250, 30, true),
('Japanese N5 Course', 'Learn Japanese from zero to N5 level', 'Lê Văn L', 'Language', 'Beginner', 50, 1800000, 4.5, 110, 25, true),

-- Music (1)
('Guitar for Beginners', 'Learn to play guitar from scratch', 'Phạm Văn M', 'Music', 'Beginner', 20, 500000, 4.8, 200, 10, true);

-- ============================================
-- Seed Lessons (for each course, at least 3 lessons)
-- ============================================

-- Course 1: Flutter Development (5 lessons)
INSERT INTO lessons (course_id, title, description, video_url, duration, "order") VALUES
(1, 'Introduction to Flutter', 'What is Flutter and why use it', 'https://example.com/flutter/1.mp4', 30, 1),
(1, 'Setting up Development Environment', 'Install Flutter SDK and IDE', 'https://example.com/flutter/2.mp4', 45, 2),
(1, 'Your First Flutter App', 'Create a simple Hello World app', 'https://example.com/flutter/3.mp4', 60, 3),
(1, 'Widgets Fundamentals', 'Understanding widgets in Flutter', 'https://example.com/flutter/4.mp4', 50, 4),
(1, 'State Management Basics', 'Managing state in Flutter apps', 'https://example.com/flutter/5.mp4', 55, 5);

-- Course 2: Advanced Python (4 lessons)
INSERT INTO lessons (course_id, title, description, video_url, duration, "order") VALUES
(2, 'Python Advanced Functions', 'Decorators, generators, and more', 'https://example.com/python/1.mp4', 40, 1),
(2, 'Object Oriented Python', 'Advanced OOP concepts', 'https://example.com/python/2.mp4', 50, 2),
(2, 'Python Multithreading', 'Concurrent programming in Python', 'https://example.com/python/3.mp4', 45, 3),
(2, 'Python Best Practices', 'Writing clean and efficient code', 'https://example.com/python/4.mp4', 35, 4);

-- Course 3: JavaScript (4 lessons)
INSERT INTO lessons (course_id, title, description, video_url, duration, "order") VALUES
(3, 'JavaScript Basics', 'Variables, types, and operators', 'https://example.com/js/1.mp4', 35, 1),
(3, 'DOM Manipulation', 'Working with the Document Object Model', 'https://example.com/js/2.mp4', 45, 2),
(3, 'Async JavaScript', 'Promises, async/await', 'https://example.com/js/3.mp4', 50, 3),
(3, 'Modern ES6+ Features', 'Arrow functions, destructuring, and more', 'https://example.com/js/4.mp4', 40, 4);

-- Course 4: React Native (3 lessons)
INSERT INTO lessons (course_id, title, description, video_url, duration, "order") VALUES
(4, 'React Native Setup', 'Setting up React Native environment', 'https://example.com/rn/1.mp4', 40, 1),
(4, 'Components and Props', 'Building reusable components', 'https://example.com/rn/2.mp4', 50, 2),
(4, 'Navigation in React Native', 'Screen navigation patterns', 'https://example.com/rn/3.mp4', 45, 3);

-- Course 5: UI/UX Design (3 lessons)
INSERT INTO lessons (course_id, title, description, video_url, duration, "order") VALUES
(5, 'Introduction to UX', 'What is user experience', 'https://example.com/ux/1.mp4', 30, 1),
(5, 'User Research Methods', 'How to understand your users', 'https://example.com/ux/2.mp4', 40, 2),
(5, 'Creating Wireframes', 'Design your first wireframe', 'https://example.com/ux/3.mp4', 45, 3);

-- Update lesson_count for courses
UPDATE courses SET lesson_count = (SELECT COUNT(*) FROM lessons WHERE lessons.course_id = courses.id) WHERE id > 0;

-- ============================================
-- Seed Enrollments (10 enrollments)
-- ============================================

INSERT INTO enrollments (student_id, course_id, enrollment_date, progress, status, certificate_issued) VALUES
(2, 1, '2024-10-01', 60, 'active', false),
(2, 3, '2024-10-15', 100, 'completed', true),
(2, 5, '2024-11-01', 30, 'active', false),
(3, 1, '2024-09-20', 80, 'active', false),
(3, 2, '2024-10-10', 45, 'active', false),
(3, 10, '2024-11-05', 20, 'active', false),
(4, 4, '2024-08-15', 100, 'completed', true),
(4, 6, '2024-09-01', 50, 'active', false),
(5, 1, '2024-11-01', 10, 'active', false),
(5, 12, '2024-10-20', 0, 'dropped', false);

-- Update student_count for courses (only count active/completed)
UPDATE courses SET student_count = (
    SELECT COUNT(*) FROM enrollments 
    WHERE enrollments.course_id = courses.id AND status IN ('active', 'completed')
) WHERE id > 0;

-- ============================================
-- Seed Lesson Progress (for some enrollments)
-- ============================================

-- Enrollment 1 (student 2, course 1): 3/5 lessons completed (60%)
INSERT INTO lesson_progress (enrollment_id, lesson_id, is_completed, completed_at) VALUES
(1, 1, true, '2024-10-02'),
(1, 2, true, '2024-10-05'),
(1, 3, true, '2024-10-10'),
(1, 4, false, NULL),
(1, 5, false, NULL);

-- Enrollment 2 (student 2, course 3): all lessons completed (100%)
INSERT INTO lesson_progress (enrollment_id, lesson_id, is_completed, completed_at) VALUES
(2, 10, true, '2024-10-20'),
(2, 11, true, '2024-10-25'),
(2, 12, true, '2024-10-30'),
(2, 13, true, '2024-11-05');

-- Verify
SELECT 'Students: ' || COUNT(*) FROM students;
SELECT 'Courses: ' || COUNT(*) FROM courses;
SELECT 'Lessons: ' || COUNT(*) FROM lessons;
SELECT 'Enrollments: ' || COUNT(*) FROM enrollments;
SELECT 'Lesson Progress: ' || COUNT(*) FROM lesson_progress;

# API Documentation - Online Learning Management System

## Base URL
```
http://localhost:8080/api
```

## Authentication
Tất cả API yêu cầu xác thực sử dụng Bearer Token trong header:
```
Authorization: Bearer <token>
```

---

## 1. Authentication

### 1.1. Đăng ký Student
**POST** `/auth/register`

**Request Body:**
```json
{
  "email": "student@example.com",
  "password": "password123",
  "full_name": "Nguyễn Văn A",
  "phone_number": "0123456789",
  "date_of_birth": "2000-01-01",
  "gender": "male"
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "Registration successful",
  "data": {
    "id": 1,
    "email": "student@example.com",
    "full_name": "Nguyễn Văn A",
    "is_active": true,
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

### 1.2. Đăng nhập
**POST** `/auth/login`

**Request Body:**
```json
{
  "email": "student@example.com",
  "password": "password123"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "student_id": "2151012345",
    "user": {
      "id": 1,
      "email": "student@example.com",
      "full_name": "Nguyễn Văn A",
      "is_active": true
    }
  }
}
```

### 1.3. Lấy thông tin Student hiện tại
**GET** `/auth/me`

**Headers:** `Authorization: Bearer <token>`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "student@example.com",
    "full_name": "Nguyễn Văn A"
  }
}
```

---

## 2. Students

### 2.1. Lấy danh sách Students (Admin only)
**GET** `/students?page=1&limit=10&search=`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "items": [...],
    "totalItems": 5,
    "currentPage": 1,
    "pageSize": 10,
    "totalPages": 1
  }
}
```

### 2.2. Lấy Student theo ID
**GET** `/students/:id`

### 2.3. Cập nhật Student
**PUT** `/students/:id`

**Request Body:**
```json
{
  "full_name": "Nguyễn Văn B",
  "phone_number": "0987654321"
}
```

---

## 3. Courses

### 3.1. Lấy danh sách Courses
**GET** `/courses`

**Query Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| page | int | Trang (default: 1) |
| limit | int | Số lượng/trang (default: 10) |
| search | string | Tìm kiếm |
| category | string | Lọc theo category |
| level | string | Lọc theo level |
| min_price | float | Giá tối thiểu |
| max_price | float | Giá tối đa |
| published_only | bool | Chỉ published (default: true) |

**Response (200):**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": 1,
        "title": "Flutter Development",
        "instructor": "John Doe",
        "category": "Programming",
        "level": "Beginner",
        "price": 1500000,
        "rating": 4.8,
        "student_count": 150,
        "lesson_count": 20
      }
    ],
    "totalItems": 12,
    "currentPage": 1,
    "pageSize": 10,
    "totalPages": 2
  }
}
```

### 3.2. Lấy Course theo ID (với Lessons)
**GET** `/courses/:id`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "title": "Flutter Development",
    "lessons": [
      {
        "id": 1,
        "title": "Introduction to Flutter",
        "duration": 30,
        "order": 1
      }
    ]
  }
}
```

### 3.3. Tạo Course (Admin only)
**POST** `/courses`

**Request Body:**
```json
{
  "title": "Flutter Development",
  "description": "Learn Flutter from scratch",
  "instructor": "John Doe",
  "category": "Programming",
  "level": "Beginner",
  "duration": 40,
  "price": 500000,
  "image_url": "https://example.com/image.jpg"
}
```

### 3.4. Thêm Lesson vào Course (Admin only)
**POST** `/courses/:id/lessons`

**Request Body:**
```json
{
  "title": "Introduction to Flutter",
  "description": "First lesson",
  "video_url": "https://example.com/video1.mp4",
  "duration": 30,
  "order": 1
}
```

---

## 4. Enrollments

### 4.1. Đăng ký Khóa học
**POST** `/enrollments`

**Request Body:**
```json
{
  "course_id": 1
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "Enrolled successfully",
  "data": {
    "id": 1,
    "student_id": 1,
    "course_id": 1,
    "progress": 0,
    "status": "active"
  }
}
```

### 4.2. Lấy Khóa học của Student
**GET** `/students/:id/enrollments?status=active&page=1&limit=10`

### 4.3. Lấy Tiến độ Khóa học
**GET** `/enrollments/:id/progress`

**Response (200):**
```json
{
  "success": true,
  "data": {
    "enrollment_id": 1,
    "course": {...},
    "progress": 45,
    "completed_lessons": 9,
    "total_lessons": 20,
    "lessons": [
      {
        "lesson_id": 1,
        "title": "Lesson 1",
        "is_completed": true,
        "completed_at": "2024-01-01T10:00:00Z"
      }
    ]
  }
}
```

### 4.4. Hoàn thành Bài học
**PUT** `/enrollments/:id/lessons/:lesson_id/complete`

**Response (200):**
```json
{
  "success": true,
  "message": "Lesson completed",
  "data": {
    "id": 1,
    "progress": 50,
    "status": "active"
  }
}
```

### 4.5. Hủy Đăng ký
**DELETE** `/enrollments/:id`

---

## Error Responses

**401 Unauthorized:**
```json
{
  "success": false,
  "message": "Invalid or expired token"
}
```

**403 Forbidden:**
```json
{
  "success": false,
  "message": "Admin access required"
}
```

**404 Not Found:**
```json
{
  "success": false,
  "message": "Resource not found"
}
```

**409 Conflict:**
```json
{
  "success": false,
  "message": "Already enrolled in this course"
}
```

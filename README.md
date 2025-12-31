# Online Learning Management API

Hệ thống quản lý học trực tuyến với RESTful API (Go + Gin + GORM) và Flutter Client App.

## 📋 Yêu cầu

- Go 1.21+
- PostgreSQL 15+
- Flutter 3.x

## 🚀 Cài đặt

### 1. Clone và cấu hình

```bash
cd backend
cp .env.example .env
# Sửa file .env với thông tin database của bạn
```

### 2. Tạo Database

```bash
# Tạo database
psql -c "CREATE DATABASE db_exam_2251172466;"

# Chạy schema
psql -d db_exam_2251172466 -f database/schema.sql

# Chạy seed data
psql -d db_exam_2251172466 -f database/seed.sql
```

### 3. Chạy Backend

```bash
cd backend
go mod tidy
go run cmd/api/main.go
```

Server chạy tại: http://localhost:8080

### 4. Chạy Flutter App

```bash
cd frontend
flutter pub get
flutter run -d chrome
```

## 🔐 Tài khoản test

| Role | Email | Password |
|------|-------|----------|
| Admin | admin@learning.com | password123 |
| Student | student1@learning.com | password123 |

## 📖 API Endpoints

### Authentication
- `POST /api/auth/register` - Đăng ký
- `POST /api/auth/login` - Đăng nhập (return student_id)
- `GET /api/auth/me` - Thông tin user hiện tại

### Students
- `GET /api/students` - Danh sách (admin)
- `GET /api/students/:id` - Chi tiết
- `PUT /api/students/:id` - Cập nhật

### Courses
- `GET /api/courses` - Danh sách (filter: search, category, level, price)
- `GET /api/courses/:id` - Chi tiết + lessons
- `POST /api/courses` - Tạo mới (admin)
- `PUT /api/courses/:id` - Cập nhật (admin)
- `DELETE /api/courses/:id` - Xóa (admin)
- `POST /api/courses/:id/lessons` - Thêm bài học (admin)

### Enrollments
- `POST /api/enrollments` - Đăng ký khóa học
- `GET /api/students/:id/enrollments` - Khóa học của student
- `GET /api/courses/:id/enrollments` - Học viên của course (admin)
- `GET /api/enrollments/:id/progress` - Tiến độ học
- `PUT /api/enrollments/:id/lessons/:lid/complete` - Hoàn thành bài học
- `DELETE /api/enrollments/:id` - Hủy đăng ký

## 📁 Cấu trúc Project

```
.
├── backend/
│   ├── cmd/api/main.go
│   ├── internal/
│   │   ├── config/
│   │   ├── models/
│   │   ├── repository/
│   │   ├── service/
│   │   ├── handler/
│   │   ├── middleware/
│   │   └── dto/
│   ├── database/
│   │   ├── schema.sql
│   │   └── seed.sql
│   └── docs/
└── frontend/
    └── lib/
        ├── models/
        ├── providers/
        ├── screens/
        └── services/
```

## 🛠 Công nghệ sử dụng

**Backend:**
- Go 1.21
- Gin Web Framework
- GORM ORM
- PostgreSQL
- JWT Authentication

**Frontend:**
- Flutter 3.x
- Provider (State Management)
- HTTP Client
# web_api_2251172466

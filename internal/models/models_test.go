package models

import (
	"testing"
)

func TestStudentIsAdmin(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Admin user", "admin@learning.com", true},
		{"Regular student", "student1@learning.com", false},
		{"Another student", "test@example.com", false},
		{"Empty email", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			student := &Student{Email: tt.email}
			if got := student.IsAdmin(); got != tt.expected {
				t.Errorf("Student.IsAdmin() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestIsValidCategory(t *testing.T) {
	tests := []struct {
		category string
		expected bool
	}{
		{"Programming", true},
		{"Design", true},
		{"Business", true},
		{"Language", true},
		{"Music", true},
		{"InvalidCategory", false},
		{"", false},
		{"programming", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.category, func(t *testing.T) {
			if got := IsValidCategory(tt.category); got != tt.expected {
				t.Errorf("IsValidCategory(%q) = %v, want %v", tt.category, got, tt.expected)
			}
		})
	}
}

func TestIsValidLevel(t *testing.T) {
	tests := []struct {
		level    string
		expected bool
	}{
		{"Beginner", true},
		{"Intermediate", true},
		{"Advanced", true},
		{"Expert", false},
		{"", false},
		{"beginner", false}, // Case sensitive
	}

	for _, tt := range tests {
		t.Run(tt.level, func(t *testing.T) {
			if got := IsValidLevel(tt.level); got != tt.expected {
				t.Errorf("IsValidLevel(%q) = %v, want %v", tt.level, got, tt.expected)
			}
		})
	}
}

func TestEnrollmentIsActive(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"Active enrollment", StatusActive, true},
		{"Completed enrollment", StatusCompleted, false},
		{"Dropped enrollment", StatusDropped, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enrollment := &Enrollment{Status: tt.status}
			if got := enrollment.IsActive(); got != tt.expected {
				t.Errorf("Enrollment.IsActive() = %v, want %v for status %s", got, tt.expected, tt.status)
			}
		})
	}
}

func TestEnrollmentIsCompleted(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		expected bool
	}{
		{"Active enrollment", StatusActive, false},
		{"Completed enrollment", StatusCompleted, true},
		{"Dropped enrollment", StatusDropped, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enrollment := &Enrollment{Status: tt.status}
			if got := enrollment.IsCompleted(); got != tt.expected {
				t.Errorf("Enrollment.IsCompleted() = %v, want %v for status %s", got, tt.expected, tt.status)
			}
		})
	}
}

func TestIsValidEnrollmentStatus(t *testing.T) {
	tests := []struct {
		status   string
		expected bool
	}{
		{"active", true},
		{"completed", true},
		{"dropped", true},
		{"pending", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			if got := IsValidEnrollmentStatus(tt.status); got != tt.expected {
				t.Errorf("IsValidEnrollmentStatus(%q) = %v, want %v", tt.status, got, tt.expected)
			}
		})
	}
}

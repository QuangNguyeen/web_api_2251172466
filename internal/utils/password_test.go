package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "password123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Hash should not be empty
	if hash == "" {
		t.Error("HashPassword returned empty hash")
	}

	// Hash should be different from password
	if hash == password {
		t.Error("Hash should not equal original password")
	}

	// Hash should start with $2a$ (bcrypt)
	if len(hash) < 4 || hash[:4] != "$2a$" {
		t.Error("Hash should be bcrypt format ($2a$...)")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "password123"
	wrongPassword := "wrongpassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Correct password should match
	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for correct password")
	}

	// Wrong password should not match
	if CheckPassword(wrongPassword, hash) {
		t.Error("CheckPassword should return false for wrong password")
	}

	// Empty password should not match
	if CheckPassword("", hash) {
		t.Error("CheckPassword should return false for empty password")
	}
}

func TestHashPasswordUniqueness(t *testing.T) {
	password := "samepassword"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	// Same password should produce different hashes (due to salt)
	if hash1 == hash2 {
		t.Error("Two hashes of the same password should be different (different salts)")
	}

	// But both should still validate
	if !CheckPassword(password, hash1) || !CheckPassword(password, hash2) {
		t.Error("Both hashes should validate the original password")
	}
}

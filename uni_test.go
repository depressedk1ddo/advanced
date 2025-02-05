package main

import (
	"testing"
)

// Test password hashing and verification
func TestPasswordHashing(t *testing.T) {
	password := "password123"
	hashedPassword, err := hashPassword(password)
	if err != nil {
		t.Fatalf("Error hashing password: %v", err)
	}

	if !checkPassword(hashedPassword, password) {
		t.Errorf("Password check failed: hashed password does not match")
	}
}

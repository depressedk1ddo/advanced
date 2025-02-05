package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginAPI(t *testing.T) {
	// Create a test user in the database
	db.Exec("DELETE FROM users WHERE email = 'test@example.com'")
	db.Exec("INSERT INTO users (name, email, password, role) VALUES ('Test User', 'test@example.com', '$2a$10$kyszjazUhUrAeIbfkRpLhOOP1R4zQSdqsp.2RBnwDueTqKCV7lYR2', 'user')")

	// Prepare login request
	loginData := map[string]string{
		"email":    "test@example.com",
		"password": "password123", // Must match the hashed password above
	}
	body, _ := json.Marshal(loginData)

	req, err := http.NewRequest("POST", "/api/login", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	recorder := httptest.NewRecorder()
	http.HandlerFunc(loginHandler).ServeHTTP(recorder, req)

	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, status)
	}
}

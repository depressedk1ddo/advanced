package main

import (
	"testing"
	"time"

	"github.com/tebeka/selenium"
)

func TestLoginPage(t *testing.T) {
	const seleniumURL = "http://localhost:4444/wd/hub"

	// Start Selenium WebDriver
	wd, err := selenium.NewRemote(selenium.Capabilities{"browserName": "chrome"}, seleniumURL)
	if err != nil {
		t.Fatal(err)
	}
	defer wd.Quit()

	// Open login page
	wd.Get("http://localhost:8080/login.html")

	// Find and fill the email input
	emailInput, err := wd.FindElement(selenium.ByID, "email")
	if err != nil {
		t.Fatal("Email input not found")
	}
	emailInput.SendKeys("test@example.com")

	// Find and fill the password input
	passwordInput, err := wd.FindElement(selenium.ByID, "password")
	if err != nil {
		t.Fatal("Password input not found")
	}
	passwordInput.SendKeys("password123")

	// Click login button
	loginButton, err := wd.FindElement(selenium.ByTagName, "button")
	if err != nil {
		t.Fatal("Login button not found")
	}
	loginButton.Click()

	// Wait for redirection
	time.Sleep(2 * time.Second)

	// Check if redirected to index.html
	currentURL, _ := wd.CurrentURL()
	if currentURL != "http://localhost:8080/index.html" {
		t.Errorf("Login failed: expected redirect to index.html, got %s", currentURL)
	}
}

package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Unit tests
// RegisterHandler simulates registering a user.
func RegisterHandler(name, email, password string) (string, error) {
	return "User registered successfully", nil
}

// TestRegisterHandler verifies the functionality of RegisterHandler.
func TestRegisterHandler(t *testing.T) {
	name := "John"
	email := "john@example.com"
	password := "password"

	expectedResponse := "User registered successfully"
	var expectedError error

	response, err := RegisterHandler(name, email, password)
	if err != expectedError || response != expectedResponse {
		t.Errorf("Unexpected result. Expected: %v, %v; Got: %v, %v", expectedResponse, expectedError, response, err)
	} else {
		t.Logf("Expected response: %v, Expected error: %v", expectedResponse, expectedError)
		t.Logf("Actual response: %v, Actual error: %v", response, err)
	}
}

// Integration tests
// GetUser simulates retrieving user details.
func GetUser(userID string) (string, error) {
	return "User details retrieved successfully", nil
}

// TestIntegrationGetUser verifies the integration of GetUser function.
func TestIntegrationGetUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userID")
		response, err := GetUser(userID)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer ts.Close()

	res, err := http.Get(ts.URL + "/getUser?userID=123")
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, res.StatusCode)
	} else {
		t.Logf("Expected status code: %d, Actual status code: %d", http.StatusOK, res.StatusCode)
	}
}

// End-to-End tests
// TestEndToEndLogin verifies the end-to-end login functionality.
func TestEndToEndLogin(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if r.URL.Path != "/login" {
			http.Error(w, "Not Found", http.StatusNotFound)
			return
		}

		email := r.FormValue("email")
		password := r.FormValue("password")

		if email == "test@example.com" && password == "password" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
		} else {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
	}))
	defer ts.Close()

	client := &http.Client{}

	payload := strings.NewReader("email=test@example.com&password=password")
	req, err := http.NewRequest("POST", ts.URL+"/login", payload)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		t.Fatalf("Error sending request: %v", err)
	}

	if res.StatusCode != http.StatusMethodNotAllowed {
		t.Errorf("Expected status code %d, got %d", http.StatusMethodNotAllowed, res.StatusCode)
	} else {
		t.Logf("Expected status code: %d, Actual status code: %d", http.StatusMethodNotAllowed, res.StatusCode)
	}

	location, err := res.Location()
	if err != nil {
		if err == http.ErrNoLocation {
			t.Logf("No redirect location found, as expected.")
		} else {
			t.Fatalf("Error getting redirect location: %v", err)
		}
	} else {
		if location.Path != "/" {
			t.Errorf("Expected redirect location '/', got %s", location.Path)
		} else {
			t.Logf("Expected redirect location: /, Actual redirect location: %s", location.Path)
		}
	}
}

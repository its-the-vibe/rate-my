package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

// TestHandleRate_ValidRating tests that a valid POST request returns 200 and a success status.
func TestHandleRate_ValidRating(t *testing.T) {
	body, err := json.Marshal(Rating{Event: "tube journey", Rating: 4})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/rate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleRate(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if result["status"] != "success" {
		t.Errorf("expected status=success, got %q", result["status"])
	}
}

// TestHandleRate_RatingBoundaries tests ratings at and beyond the valid 1–5 range.
func TestHandleRate_RatingBoundaries(t *testing.T) {
	tests := []struct {
		name           string
		rating         int
		expectedStatus int
	}{
		{"minimum valid", 1, http.StatusOK},
		{"maximum valid", 5, http.StatusOK},
		{"below minimum", 0, http.StatusBadRequest},
		{"above maximum", 6, http.StatusBadRequest},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body, err := json.Marshal(Rating{Event: "concert", Rating: tc.rating})
			if err != nil {
				t.Fatalf("failed to marshal request body: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, "/rate", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			handleRate(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("rating %d: expected %d, got %d", tc.rating, tc.expectedStatus, w.Code)
			}
		})
	}
}

// TestHandleRate_InvalidJSON tests that malformed JSON returns 400.
func TestHandleRate_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/rate", bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleRate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

// TestHandleRate_MethodNotAllowed tests that non-POST methods return 405.
func TestHandleRate_MethodNotAllowed(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodPut, http.MethodDelete} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/rate", nil)
			w := httptest.NewRecorder()

			handleRate(w, req)

			if w.Code != http.StatusMethodNotAllowed {
				t.Errorf("method %s: expected 405, got %d", method, w.Code)
			}
		})
	}
}

// TestHandleRate_TimestampDefaulted tests that an empty timestamp is filled in automatically.
func TestHandleRate_TimestampDefaulted(t *testing.T) {
	body, err := json.Marshal(Rating{Event: "bus ride", Rating: 3})
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/rate", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handleRate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

// TestAppendRatingToFile tests that appendRatingToFile writes to a file correctly.
func TestAppendRatingToFile(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "ratings.log")

	// Temporarily override the working directory so appendRatingToFile writes to our temp dir.
	original, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("could not change directory: %v", err)
	}
	defer os.Chdir(original)

	line := `{"event":"test","rating":5}`
	if err := appendRatingToFile(line); err != nil {
		t.Fatalf("appendRatingToFile returned error: %v", err)
	}

	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("could not read log file: %v", err)
	}
	if string(data) != line+"\n" {
		t.Errorf("unexpected file content: %q", string(data))
	}
}

// TestLoadRedisConfig tests parsing a valid Redis YAML config.
func TestLoadRedisConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "redis.yml")

	yaml := "host: localhost:6379\nlist: ratings\n"
	if err := os.WriteFile(cfgPath, []byte(yaml), 0644); err != nil {
		t.Fatalf("could not write config: %v", err)
	}

	cfg, err := loadRedisConfig(cfgPath)
	if err != nil {
		t.Fatalf("loadRedisConfig returned error: %v", err)
	}
	if cfg.Host != "localhost:6379" {
		t.Errorf("expected host=localhost:6379, got %q", cfg.Host)
	}
	if cfg.List != "ratings" {
		t.Errorf("expected list=ratings, got %q", cfg.List)
	}
}

// TestLoadRedisConfig_Missing tests that a missing file returns an error.
func TestLoadRedisConfig_Missing(t *testing.T) {
	_, err := loadRedisConfig("/nonexistent/path/redis.yml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Rating represents a user's rating submission
type Rating struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	Rating    int    `json:"rating"`
}

var logFile *os.File

func main() {
	// Open or create the ratings log file for appending
	var err error
	logFile, err = os.OpenFile("ratings.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer logFile.Close()

	// Serve static files from the static directory
	http.Handle("/", http.FileServer(http.Dir("static")))

	// Handle rating submissions
	http.HandleFunc("/rate", handleRate)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleRate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var rating Rating
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&rating); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate rating is between 1 and 5
	if rating.Rating < 1 || rating.Rating > 5 {
		http.Error(w, "Rating must be between 1 and 5", http.StatusBadRequest)
		return
	}

	// If no timestamp provided, use current time
	if rating.Timestamp == "" {
		rating.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}

	// Log to stdout (primary logging)
	logLine, _ := json.Marshal(rating)
	fmt.Println(string(logLine))

	// Append to file (secondary logging - non-fatal if it fails)
	if _, err := logFile.WriteString(string(logLine) + "\n"); err != nil {
		log.Printf("Failed to write to log file: %v", err)
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

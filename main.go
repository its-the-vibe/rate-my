package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
)

// Rating represents a user's rating submission
type Rating struct {
	Timestamp string `json:"timestamp"`
	Event     string `json:"event"`
	Rating    int    `json:"rating"`
	Comment   string `json:"comment"`
}

// redisConfig holds the Redis connection configuration loaded from config/redis.yml
type redisConfig struct {
	Host string `yaml:"host"`
	List string `yaml:"list"`
}

// redisClient is the global Redis client; nil if Redis is not configured
var redisClient *redis.Client

// redisListName is the Redis list key to push log lines to
var redisListName string

func main() {
	// Load .env file if present (best-effort; errors are non-fatal)
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Could not load .env file: %v", err)
		}
	}

	// Load Redis config if the config file exists
	if cfg, err := loadRedisConfig("config/redis.yml"); err == nil {
		password := os.Getenv("REDIS_PASSWORD")
		redisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.Host,
			Password: password,
		})
		redisListName = cfg.List
		// Verify connectivity; log a warning but continue if the ping fails
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			log.Printf("Redis ping failed (host=%s): %v", cfg.Host, err)
		} else {
			log.Printf("Redis configured: host=%s list=%s", cfg.Host, cfg.List)
		}
	} else if !os.IsNotExist(err) {
		log.Printf("Could not load Redis config: %v", err)
	}

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
	if err := appendRatingToFile(string(logLine)); err != nil {
		log.Printf("Failed to write to log file: %v", err)
	}

	// Push to Redis list (non-fatal if it fails or Redis is not configured)
	if redisClient != nil {
		if err := rpushLogLine(r.Context(), string(logLine)); err != nil {
			log.Printf("Failed to push log line to Redis: %v", err)
		}
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// appendRatingToFile appends a single log line to ratings.log, opening and closing the file each time
func appendRatingToFile(logLine string) error {
	f, err := os.OpenFile("ratings.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(logLine + "\n")

	return err
}

// loadRedisConfig reads and parses the YAML Redis configuration file at the given path.
func loadRedisConfig(path string) (*redisConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg redisConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// rpushLogLine pushes a log line to the configured Redis list using RPUSH.
func rpushLogLine(ctx context.Context, logLine string) error {
	return redisClient.RPush(ctx, redisListName, logLine).Err()
}

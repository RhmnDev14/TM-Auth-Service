package helper

import (
	"auth-service/internal/dto"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// GetEnvString mengambil environment variable, jika tidak ada, log warning dan kembalikan string kosong.
func GetEnvString(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("⚠️  Warning: Environment variable %s not set\n", key)
	}
	return value
}

// GetEnvInt mengambil environment variable int, jika tidak valid atau kosong, log warning dan kembalikan 0.
func GetEnvInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("⚠️  Warning: Environment variable %s not set\n", key)
		return 0
	}
	i, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("⚠️  Warning: Environment variable %s must be an integer\n", key)
		return 0
	}
	return i
}

// GetEnvBool mengambil environment variable boolean, jika tidak valid atau kosong, log warning dan kembalikan false.
func GetEnvBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("⚠️  Warning: Environment variable %s not set\n", key)
		return false
	}
	b, err := strconv.ParseBool(value)
	if err != nil {
		log.Printf("⚠️  Warning: Environment variable %s must be a boolean (true/false)\n", key)
		return false
	}
	return b
}

// GetEnvDuration mengambil environment variable int dalam menit, dikonversi ke time.Duration.
// Jika tidak valid atau kosong, log warning dan kembalikan 0.
func GetEnvDuration(key string) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("⚠️  Warning: Environment variable %s not set\n", key)
		return 0
	}
	i, err := strconv.Atoi(value)
	if err != nil || i <= 0 {
		log.Printf("⚠️  Warning: Environment variable %s must be a positive integer\n", key)
		return 0
	}
	return time.Duration(i) * time.Minute
}

// WriteJSON adalah helper untuk menulis respons JSON
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if msg, ok := data.(string); ok {
		json.NewEncoder(w).Encode(dto.Response{
			Status:  status,
			Message: msg,
		})
		return
	}

	json.NewEncoder(w).Encode(data)
}

func ErrorHandle(param error) error {
	return fmt.Errorf("ERROR : %v", param)
}

func ParseUint(s string) (uint, error) {
	val, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(val), nil
}

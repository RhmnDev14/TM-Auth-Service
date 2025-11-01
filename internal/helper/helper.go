package helper

import (
	"auth-service/internal/dto"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

// GetEnvRequired mengambil nilai dari environment variable. Jika kosong, ia akan panic.
func GetEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("FATAL: Environment variable %s is required and not set", key))
	}
	return value
}

// GetEnvInt mengambil nilai int dari environment variable.
func GetEnvInt(key string) int {
	value := GetEnvRequired(key)
	i, err := strconv.Atoi(value)
	if err != nil {
		panic(fmt.Sprintf("FATAL: Environment variable %s must be an integer: %v", key, err))
	}
	return i
}

// GetEnvBool mengambil nilai boolean dari environment variable.
func GetEnvBool(key string) bool {
	value := GetEnvRequired(key)
	b, err := strconv.ParseBool(value)
	if err != nil {
		panic(fmt.Sprintf("FATAL: Environment variable %s must be a boolean (true/false): %v", key, err))
	}
	return b
}

// GetEnvDuration mengambil nilai duration dari environment variable.
func GetEnvDuration(key string) time.Duration {
	value := GetEnvInt(key)
	if value <= 0 {
		panic(fmt.Sprintf("FATAL: Environment variable %s must be a positive integer", key))
	}
	// Konversi nilai integer (menit) ke time.Duration
	return time.Duration(value) * time.Minute
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

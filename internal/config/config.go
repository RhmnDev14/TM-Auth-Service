package config

import (
	"auth-service/internal/helper"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	Name        string
	LogMode     bool
	MaxIdle     int
	MaxOpen     int
	MaxLife     int
	MaxIdleTime int
}

type TokenConfig struct {
	IssuerName              string
	JwtSignatureKy          []byte
	JwtExpiresTime          time.Duration
	RefreshTokenExpiresTime time.Duration
	JwtSigningMethod        *jwt.SigningMethodHMAC
}

type PortConfig struct {
	ServerPort string
}

type Config struct {
	DBConfig
	TokenConfig
	PortConfig
}

func NewConfig() *Config {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		panic("⚠️  Warning: .env file not found, using system environment variables")
	}

	var cfg Config

	// Load DBConfig
	cfg.DBConfig = DBConfig{
		Host:        helper.GetEnvRequired("DB_HOST"),
		Port:        helper.GetEnvRequired("DB_PORT"),
		User:        helper.GetEnvRequired("DB_USER"),
		Password:    helper.GetEnvRequired("DB_PASSWORD"),
		Name:        helper.GetEnvRequired("DB_NAME"),
		LogMode:     helper.GetEnvBool("DB_LOG_MODE"),
		MaxIdle:     helper.GetEnvInt("DB_MAX_IDLE_CONNS"),
		MaxOpen:     helper.GetEnvInt("DB_MAX_OPEN_CONNS"),
		MaxLife:     helper.GetEnvInt("DB_MAX_LIFE_TIME"),
		MaxIdleTime: helper.GetEnvInt("DB_MAX_IDLE_TIME"),
	}

	// Load PortConfig
	cfg.PortConfig = PortConfig{
		ServerPort: helper.GetEnvRequired("SERVER_PORT"),
	}

	//  Load TokenConfig
	jwtSecret := helper.GetEnvRequired("TOKEN_SECRET")

	cfg.TokenConfig = TokenConfig{
		IssuerName:     helper.GetEnvRequired("TOKEN_ISSUE"),
		JwtSignatureKy: []byte(jwtSecret),
		// TOKEN_EXPIRE & REFRESH_TOKEN_EXPIRE nilainya integer dalam MENIT
		JwtExpiresTime:          helper.GetEnvDuration("TOKEN_EXPIRE"),
		RefreshTokenExpiresTime: helper.GetEnvDuration("REFRESH_TOKEN_EXPIRE"),
		JwtSigningMethod:        jwt.SigningMethodHS256,
	}

	// Validasi
	if err := cfg.validate(); err != nil {
		panic(fmt.Sprintf("FATAL: Configuration validation failed: %v", err))
	}

	return &cfg
}

func (c *Config) validate() error {
	if len(c.TokenConfig.JwtSignatureKy) == 0 {
		return fmt.Errorf("JWT signature key is required")
	}

	return nil
}

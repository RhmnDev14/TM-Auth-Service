package config

import (
	"auth-service/internal/helper"
	"log"
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
		log.Println("⚠️  Warning: .env file not found, using system environment variables")
	}

	var cfg Config

	// Load DBConfig
	cfg.DBConfig = DBConfig{
		Host:        helper.GetEnvString("DB_HOST"),
		Port:        helper.GetEnvString("DB_PORT"),
		User:        helper.GetEnvString("DB_USER"),
		Password:    helper.GetEnvString("DB_PASSWORD"),
		Name:        helper.GetEnvString("DB_NAME"),
		LogMode:     helper.GetEnvBool("DB_LOG_MODE"),
		MaxIdle:     helper.GetEnvInt("DB_MAX_IDLE_CONNS"),
		MaxOpen:     helper.GetEnvInt("DB_MAX_OPEN_CONNS"),
		MaxLife:     helper.GetEnvInt("DB_MAX_LIFE_TIME"),
		MaxIdleTime: helper.GetEnvInt("DB_MAX_IDLE_TIME"),
	}

	// Load PortConfig
	cfg.PortConfig = PortConfig{
		ServerPort: helper.GetEnvString("SERVER_PORT"),
	}

	//  Load TokenConfig
	jwtSecret := helper.GetEnvString("TOKEN_SECRET")

	cfg.TokenConfig = TokenConfig{
		IssuerName:     helper.GetEnvString("TOKEN_ISSUE"),
		JwtSignatureKy: []byte(jwtSecret),
		// TOKEN_EXPIRE & REFRESH_TOKEN_EXPIRE nilainya integer dalam MENIT
		JwtExpiresTime:          helper.GetEnvDuration("TOKEN_EXPIRE"),
		RefreshTokenExpiresTime: helper.GetEnvDuration("REFRESH_TOKEN_EXPIRE"),
		JwtSigningMethod:        jwt.SigningMethodHS256,
	}

	// Validasi
	cfg.validate()
	return &cfg
}

func (c *Config) validate() {
	if len(c.TokenConfig.JwtSignatureKy) == 0 {
		log.Println("⚠️  Warning: JWT signature key is empty")
	}
}

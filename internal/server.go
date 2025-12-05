package internal

import (
	"auth-service/internal/config"
	"auth-service/internal/entity"
	"auth-service/internal/handler"
	"auth-service/internal/helper"
	"auth-service/internal/infrastructure"
	"auth-service/internal/repository"
	"auth-service/internal/usecase"
	"fmt"
	"log"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// -----------------------------------------------------------------

type Server struct {
	authUc usecase.AuthUc
	host   string
	router *http.ServeMux
}

// Middleware untuk CORS
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Izinkan domain frontend kamu di sini (misal: http://localhost:5173 atau domain production)
		w.Header().Set("Access-Control-Allow-Origin", "*") // gunakan "*" hanya untuk dev
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight request (OPTIONS)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) initRoute() {
	// Inisialisasi Handler Auth
	authMux := http.NewServeMux()
	authHandler := handler.NewAuthHandler(s.authUc, authMux)

	authHandler.SetupRoutes()

	// Bungkus authMux dengan middleware CORS
	s.router.Handle(helper.ApiGrup+"/", corsMiddleware(http.StripPrefix(helper.ApiGrup, authMux)))

	log.Println("✅ROUTES SETUP COMPLETE ON PREFIX : ", helper.ApiGrup)
}

func (s *Server) Run() {
	s.initRoute()
	log.Printf("🚀AUTH SERVICE STARTING ON HOST %s", s.host)

	if err := http.ListenAndServe(s.host, s.router); err != nil {
		panic(fmt.Errorf("server not running on host %s, because of error %v", s.host, err))
	}
}

func NewServer() *Server {
	// Load Konfigurasi
	cfg := config.NewConfig()

	// Koneksi Database (Menggunakan GORM dan driver Postgres)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBConfig.Host, cfg.DBConfig.User, cfg.DBConfig.Password, cfg.DBConfig.Name, cfg.DBConfig.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(fmt.Errorf("failed to connect database : %v", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(fmt.Errorf("failed to get underlying *sql.DB: %v", err))
	}

	// Set konfigurasi koneksi pool
	sqlDB.SetMaxIdleConns(cfg.DBConfig.MaxIdle)
	sqlDB.SetMaxOpenConns(cfg.DBConfig.MaxOpen)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.DBConfig.MaxLife) * time.Minute)
	sqlDB.SetConnMaxIdleTime(time.Duration(cfg.DBConfig.MaxIdleTime) * time.Minute)

	// Cek koneksi
	if err = sqlDB.Ping(); err != nil {
		panic(fmt.Errorf("failed to ping database: %v", err))
	}

	// Migrasi Database
	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		panic(fmt.Errorf("GORM AutoMigrate failed: %v", err))
	}
	log.Println("✅CONNECTING DATABASE SUCCES")

	// Inisialisasi Server
	host := fmt.Sprintf(":%s", cfg.ServerPort)
	router := http.NewServeMux()

	// Dependency Injection (Konstruksi Service Layer)
	authRepo := repository.NewAuthRepo(db)

	jwt := infrastructure.NewJWTService(*cfg)

	authUc := usecase.NewAuthUc(authRepo, jwt)

	return &Server{
		authUc: authUc,
		host:   host,
		router: router,
	}
}

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ddcringe/SMT_showdown/internal/auth"
	"github.com/ddcringe/SMT_showdown/internal/repository"
	"github.com/ddcringe/SMT_showdown/pkg/database"
	"github.com/ddcringe/SMT_showdown/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg := loadConfig()

	// Инициализация БД
	db, err := database.NewPostgresDB(cfg.DB)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Инициализация зависимостей
	userRepo := repository.NewUserRepository(db)
	authService := auth.NewAuthService(userRepo)
	jwtUtil := jwt.NewTokenUtil(cfg.JWT.Secret, cfg.JWT.TTL)
	authHandler := auth.NewHandler(authService, jwtUtil)

	// Настройка роутера
	router := gin.Default()

	// Public routes
	router.POST("/api/register", authHandler.Register)
	router.POST("/api/login", authHandler.Login)

	// Protected routes
	authGroup := router.Group("/api")
	authGroup.Use(auth.AuthMiddleware(jwtUtil))
	{
		authGroup.GET("/profile", func(c *gin.Context) {
			userID := c.MustGet("userID").(int)
			c.JSON(http.StatusOK, gin.H{"user_id": userID})
		})
	}

	// Запуск сервера
	srv := &http.Server{
		Addr:    ":" + cfg.Server.Port,
		Handler: router,
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}

type Config struct {
	Server struct {
		Port string
	}
	DB struct {
		Host     string
		Port     string
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	JWT struct {
		Secret string
		TTL    time.Duration
	}
}

func loadConfig() *Config {
	var cfg Config

	// Сервер
	cfg.Server.Port = os.Getenv("SERVER_PORT")
	if cfg.Server.Port == "" {
		cfg.Server.Port = "8080"
	}

	// База данных
	cfg.DB.Host = os.Getenv("DB_HOST")
	cfg.DB.Port = os.Getenv("DB_PORT")
	cfg.DB.User = os.Getenv("DB_USER")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")
	cfg.DB.Name = os.Getenv("DB_NAME")
	cfg.DB.SSLMode = os.Getenv("DB_SSLMODE")
	if cfg.DB.SSLMode == "" {
		cfg.DB.SSLMode = "disable"
	}

	// JWT
	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	if cfg.JWT.Secret == "" {
		log.Fatal("JWT_SECRET must be set")
	}
	cfg.JWT.TTL = 24 * time.Hour

	return &cfg
}

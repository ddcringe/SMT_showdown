package main

import (
	"SMT_showdown/internal/handlers"
	"SMT_showdown/internal/repository"
	"SMT_showdown/pkg/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка .env файла
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	// Инициализация подключения к БД
	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Инициализация репозитория
	repo := repository.NewRepository(db)

	// Создание Gin роутера
	r := gin.Default()

	// Настройка CORS
	r.Use(utils.CORSMiddleware())

	// Регистрация обработчиков
	handlers.RegisterAuthHandlers(r, repo)
	handlers.RegisterUserHandlers(r, repo)
	handlers.RegisterTeamHandlers(r, repo)
	handlers.RegisterDemonHandlers(r, repo)
	handlers.RegisterBattleHandlers(r, repo)

	// Запуск сервера
	port := utils.GetEnv("PORT", "8080")
	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

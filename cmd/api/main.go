package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dinosaur1258/GolangFramework/internal/handler"
	"github.com/dinosaur1258/GolangFramework/internal/repository/postgres"
	"github.com/dinosaur1258/GolangFramework/internal/router"
	"github.com/dinosaur1258/GolangFramework/internal/service"
	"github.com/dinosaur1258/GolangFramework/internal/usecase"
	"github.com/dinosaur1258/GolangFramework/pkg/config"
	"github.com/dinosaur1258/GolangFramework/pkg/database"
)

func main() {
	// 根據環境選擇配置檔案
	configPath := "config/config.yaml"
	if os.Getenv("DOCKER_ENV") == "true" {
		configPath = "config/config.docker.yaml"
	}

	// 載入配置
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 建立資料庫連線
	dbConfig := database.Config{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		DBName:   cfg.Database.DBName,
		SSLMode:  cfg.Database.SSLMode,
	}

	db, err := database.NewPostgresDB(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("✅ Database connected successfully!")

	// 初始化 Services
	jwtService := service.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// 依賴注入：Repository -> UseCase -> Handler
	userRepo := postgres.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := handler.NewUserHandler(userUseCase)
	authHandler := handler.NewAuthHandler(userUseCase, jwtService)

	// 設定路由
	r := router.SetupRouter(userHandler, authHandler, jwtService)

	// 啟動伺服器
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("🚀 Server is running on %s", addr)

	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

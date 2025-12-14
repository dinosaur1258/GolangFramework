package main

// @title           Golang Clean Architecture API
// @version         1.0
// @description     這是一個使用 Gin 框架和乾淨架構的 RESTful API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 輸入 "Bearer {token}" 來進行認證

import (
	"fmt"
	"log"
	"os"

	_ "github.com/dinosaur1258/GolangFramework/docs"
	"github.com/dinosaur1258/GolangFramework/internal/handler"
	"github.com/dinosaur1258/GolangFramework/internal/repository/postgres"
	"github.com/dinosaur1258/GolangFramework/internal/router"
	"github.com/dinosaur1258/GolangFramework/internal/service"
	"github.com/dinosaur1258/GolangFramework/internal/usecase"
	"github.com/dinosaur1258/GolangFramework/pkg/config"
	"github.com/dinosaur1258/GolangFramework/pkg/database"
	"github.com/dinosaur1258/GolangFramework/pkg/logger"
	"go.uber.org/zap"
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

	// 初始化日誌系統
	env := "development" // 或從 cfg.Server.Mode 讀取
	if err := logger.InitLogger(env); err != nil {
		log.Fatal("Failed to initialize logger:", err)
	}
	defer logger.Sync()

	logger.Info("Application starting", zap.String("env", env))

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

	logger.Info("Database connected successfully")

	// 初始化 Services
	jwtService := service.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// 依賴注入：Repository -> UseCase -> Handler
	userRepo := postgres.NewUserRepository(db)

	// 建立 UseCase
	// Auth 使用 UserRepository 因為認證需要查詢用戶資料
	authUseCase := usecase.NewAuthUseCase(userRepo) // ← 新增這行
	userUseCase := usecase.NewUserUseCase(userRepo)

	// 建立 Handler
	authHandler := handler.NewAuthHandler(authUseCase, jwtService) // ← 修改這行
	userHandler := handler.NewUserHandler(userUseCase)

	// 設定路由
	r := router.SetupRouter(userHandler, authHandler, jwtService)

	// 啟動伺服器
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	logger.Info("Server starting", zap.String("addr", addr))

	if err := r.Run(addr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

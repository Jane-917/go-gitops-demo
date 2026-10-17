package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Jane-917/go-gitops-demo/config"


	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Task struct {
	gorm.Model
	Title string `json:"title" binding:"required"`
	Done  bool   `json:"done"`
}

var (
	DB     *gorm.DB
	Logger *zap.Logger
)

func main() {
	var err error

	// 1. 初始化 Zap Logger
	Logger, err = zap.NewProduction()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}
	defer Logger.Sync()

	// 2. 載入 Viper 組態設定
	cfg, err := config.LoadConfig()
	if err != nil {
		Logger.Fatal("Failed to load configuration", zap.Error(err))
	}
	Logger.Info("Configuration loaded successfully")

	// 3. 使用 Viper 設定組裝 PostgreSQL DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.DBTimeZone,
	)

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		Logger.Fatal("Failed to connect database", zap.Error(err))
	}
	Logger.Info("Database connected successfully!")

	// 自動遷移資料表
	DB.AutoMigrate(&Task{})

	// 4. 初始化 Gin
	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "db": "connected"})
	})

	r.GET("/api/v1/tasks", func(c *gin.Context) {
		var tasks []Task
		DB.Find(&tasks)
		c.JSON(http.StatusOK, gin.H{"code": 200, "data": tasks})
	})

	r.POST("/api/v1/tasks", func(c *gin.Context) {
		var input Task
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		task := Task{Title: input.Title, Done: false}
		DB.Create(&task)
		c.JSON(http.StatusCreated, gin.H{"code": 201, "data": task})
	})

	// 5. 使用動態監聽 Port 啟動服務
	serverAddr := fmt.Sprintf(":%s", cfg.Port)
	Logger.Info("Server running", zap.String("address", serverAddr))
	r.Run(serverAddr)
}
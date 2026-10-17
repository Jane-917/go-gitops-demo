package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 結構體映射我們的所有設定
type Config struct {
	Port       string `mapstructure:"PORT"`
	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`
	DBTimeZone string `mapstructure:"DB_TIMEZONE"`
}

// LoadConfig 讀取設定檔與環境變數
func LoadConfig() (config Config, err error) {
	// 設定讀取路徑與檔名
	viper.AddConfigPath(".")
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// 允許 Viper 自動讀取系統環境變數（Overriding）
	viper.AutomaticEnv()

	// 嘗試讀取設定檔
	err = viper.ReadInConfig()
	if err != nil {
		log.Printf("Warning: Could not find .env file, relying on environment variables: %v", err)
	}

	// 將設定 Unmarshal 到 Config struct 中
	err = viper.Unmarshal(&config)
	return
}
package storage

import (
	"fmt"
	"log"
	"velox-raptor-saham/include/models"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB(cfg *viper.Viper) error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.GetString("database.host"),
		cfg.GetString("database.port"),
		cfg.GetString("database.user"),
		cfg.GetString("database.password"),
		cfg.GetString("database.name"),
		cfg.GetString("database.sslmode"),
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return err
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxOpenConns(cfg.GetInt("database.max_open_conn"))
	sqlDB.SetMaxIdleConns(cfg.GetInt("database.max_idle_conn"))

	// Auto Migrate
	err = DB.AutoMigrate(&models.User{})
	if err != nil {
		log.Printf("Warning: Migration failed: %v", err)
	}

	return nil
}

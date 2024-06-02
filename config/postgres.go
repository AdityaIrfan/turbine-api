package config

import (
	"errors"
	"os"

	"github.com/phuslu/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres() *gorm.DB {
	dsn := "host=localhost user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(errors.New("=== ERROR DATABASE CONNECTION : " + err.Error()))
		os.Exit(1)
	}

	log.Info().Msg("=== DATABASE CONNECTION SUCCESSFULLY ===")
	return db
}

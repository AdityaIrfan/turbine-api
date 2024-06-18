package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/phuslu/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres() *gorm.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	if _, err := strconv.Atoi(port); err != nil {
		log.Error().Err(errors.New("DB_PORT IS NOT NUMBER, CHECK ENV")).Msg("")
		os.Exit(1)
	}
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_DATABASE")
	timezone := os.Getenv("DB_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", host, user, password, databaseName, port, timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(errors.New("=== ERROR DATABASE CONNECTION : " + err.Error())).Msg("")
		os.Exit(1)
	}

	log.Info().Msg("=== DATABASE CONNECTION SUCCESSFULLY ===")
	return db
}

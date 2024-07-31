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
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	if _, err := strconv.Atoi(port); err != nil {
		log.Error().Err(errors.New("POSTGRES_PORT IS NOT NUMBER, CHECK ENV")).Msg("")
		os.Exit(1)
	}
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	databaseName := os.Getenv("POSTGRES_DATABASE")
	timezone := os.Getenv("POSTGRES_TIMEZONE")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s TimeZone=%s", host, user, password, databaseName, port, timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(errors.New("=== ERROR DATABASE CONNECTION : " + err.Error())).Msg("")
		os.Exit(1)
	}

	log.Info().Msg("=== DATABASE CONNECTION SUCCESSFULLY ===")
	return db
}

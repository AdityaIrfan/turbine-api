package config

import (
	"errors"
	"os"

	"github.com/phuslu/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitPostgres() *gorm.DB {
	dsn := "host=103.59.94.19 user=postgres password=Jbsd8she2j3neoads231j@*7jn dbname=turbine-app port=5432 TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Error().Err(errors.New("=== ERROR DATABASE CONNECTION : " + err.Error())).Msg("")
		os.Exit(1)
	}

	log.Info().Msg("=== DATABASE CONNECTION SUCCESSFULLY ===")
	return db
}

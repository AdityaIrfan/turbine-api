package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"turbine-api/config"
	"turbine-api/routes"

	"github.com/phuslu/log"
)

func main() {
	port := 8081
	defer func() {
		if err := recover(); err != nil {
			log.Error().Err(fmt.Errorf("FAILED RECOVER : %v", err)).Msg("")
			os.Exit(1)
		}
	}()

	config.InitLogger()
	// postgres := config.InitPostgres()
	// redis := config.InitRedis()

	// apis := routes.NewApi().Init(postgres, redis)
	apis := routes.NewApi().Init(nil, nil)

	log.Info().Msg(fmt.Sprintf("This server is running on port %d", port))

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), apis); err != nil {
		log.Error().Err(errors.New("FAILED TO RUN SERVER : " + err.Error())).Msg("")
		os.Exit(1)
	}
}

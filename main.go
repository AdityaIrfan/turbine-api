package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"pln/AdityaIrfan/turbine-api/config"
	"pln/AdityaIrfan/turbine-api/helpers"
	"pln/AdityaIrfan/turbine-api/routes"

	"github.com/joho/godotenv"
	"github.com/phuslu/log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Error().Err(fmt.Errorf("FAILED LOAD ENV : %v", err)).Msg("")
		os.Exit(1)
	}

	port, err := strconv.Atoi(os.Getenv("SERVER_PORT"))
	if err != nil {
		log.Error().Err(errors.New("SERVER_PORT IS NOT A NUMBER, CHECK ENV")).Msg("")
		os.Exit(1)
	}
	defer func() {
		if err := recover(); err != nil {
			log.Error().Err(fmt.Errorf("FAILED RECOVER : %v", err)).Msg("")
			os.Exit(1)
		}
	}()

	helpers.LoadConstData()
	config.InitLogger()
	postgres := config.InitPostgres()
	redis := config.InitRedis()

	apis := routes.NewApi().Init(postgres, redis)

	log.Info().Msg(fmt.Sprintf("This server is running on port %d", port))

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), apis); err != nil {
		log.Error().Err(errors.New("FAILED TO RUN SERVER : " + err.Error())).Msg("")
		os.Exit(1)
	}
}

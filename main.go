package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"turbine-api/config"
	"turbine-api/routes"
)

func main() {
	port := 8081
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("Panic Detected : " + fmt.Sprint(err))
		}
	}()

	postgres := config.InitPostgres()
	redis := config.InitRedis()

	apis := routes.NewApi().Init(postgres, redis)

	fmt.Println("This server is running on port", strconv.Itoa(port))

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), apis); err != nil {
		log.Println("FAILED TO RUN SERVER : " + err.Error())
	}
}

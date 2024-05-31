package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"turbine-api/routes"
)

func main() {
	port := 8081
	defer func() {
		if err := recover(); err != nil {
			log.Fatal("Panic Detected : " + fmt.Sprint(err))
		}
	}()

	fmt.Println("This server is running on port", strconv.Itoa(port))

	if err := http.ListenAndServe(fmt.Sprintf(":%v", port), routes.NewApi().Init()); err != nil {
		log.Println("FAILED TO RUN SERVER : " + err.Error())
	}
}

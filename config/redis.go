package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	// redisDb, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	// if err != nil {
	// 	log.Println("error redis db : ", err)
	// 	return
	// }

	// log.Println(os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"))

	host := os.Getenv("REDIS_HOST")
	port, err := strconv.Atoi(os.Getenv("REDIS_PORT"))
	if err != nil {
		log.Error().Err(errors.New("REDIS PORT IS NOT A NUMBER, CHECK ENV")).Msg("")
		os.Exit(1)
	}
	password := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:       fmt.Sprintf("%s:%d", host, port),
		Password:   password,
		DB:         0,
		MaxRetries: 3,
		PoolSize:   200,
	})

	// client.Do(context.Background(), "CLIENT", "SETNAME", "saas-be-profile-manager")
	// fmt.Println(client.ClientGetName(context.Background()))

	if err := client.Ping(context.Background()).Err(); err != nil {
		log.Error().Err(errors.New("REDIS PING INITIALIZATION FAILED : " + err.Error())).Msg("")
		os.Exit(1)
	}

	log.Info().Msg("=== REDIS CONNECTION SUCCESSFULLY ===")

	return client
}

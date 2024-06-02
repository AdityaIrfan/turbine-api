package config

import (
	"context"
	"fmt"

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

	client := redis.NewClient(&redis.Options{
		Addr:       "host:9231",
		Password:   "passowrd",
		DB:         0,
		MaxRetries: 3,
		PoolSize:   200,
	})

	client.Do(context.Background(), "CLIENT", "SETNAME", "saas-be-profile-manager")
	fmt.Println(client.ClientGetName(context.Background()))

	log.Info().Msg("=== REDIS CONNECTION SUCCESSFULLY ===")

	return client
}

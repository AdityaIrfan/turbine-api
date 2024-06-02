package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
)

type configRedisRepository struct {
	client *redis.Client
}

func NewConfigRedisRepository(client *redis.Client) contract.IConfigRedisRepository {
	return &configRedisRepository{
		client: client,
	}
}

func (a *configRedisRepository) SaveRootLocation(rootLocation *models.ConfigRootLocation) {
	key := "config_root-location"

	json, err := json.Marshal(rootLocation)
	if err != nil {
		log.Error().Err(errors.New("ERROR MARSHAL ROOT LOCATION : " + err.Error()))
		return
	}

	if err := a.client.Set(context.Background(), key, json, helpers.RootLocationRedisExpiration).Err(); err != nil {
		log.Error().Err(errors.New("ERROR SAVING ROOT LOCATION ON REDIS : " + err.Error()))
	}
}

func (a *configRedisRepository) GetRootLocation() (*models.ConfigRootLocation, error) {
	key := "config_root-location"

	val, err := a.client.Get(context.Background(), key).Result()
	if err != nil {
		log.Error().Err(errors.New("ERROR GETTING ROOT LOCATION ON REDIS : " + err.Error()))
		return nil, nil
	}

	if val == "" {
		log.Error().Err(fmt.Errorf("ERROR REDIS KEY %s : DATA IS EMPTY", key))
		return nil, nil
	}

	var refreshtoken *models.ConfigRootLocation
	if err := json.Unmarshal([]byte(val), &refreshtoken); err != nil {
		log.Error().Err(fmt.Errorf("ERROR UNMARSHAL ROOT LOCATION REDIS THAT HAVING KEY %s : %s", key, err.Error()))
		return nil, err
	}

	return refreshtoken, nil
}
package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	contract "turbine-api/contracts"
	"turbine-api/helpers"
	"turbine-api/models"

	"github.com/phuslu/log"
	"github.com/redis/go-redis/v9"
)

type authRedisRepository struct {
	client *redis.Client
}

func NewAuthRedisRepository(client *redis.Client) contract.IAuthRedisRepository {
	return &authRedisRepository{
		client: client,
	}
}

func (a *authRedisRepository) SaveRefreshToken(id string, refreshToken *models.RefreshTokenRedis, ttl time.Duration) {
	key := fmt.Sprintf("%s_refresh-token", id)

	value, err := json.Marshal(refreshToken)
	if err != nil {
		log.Error().Err(errors.New("ERROR MARSHAL REFRESH TOKEN : " + err.Error())).Msg("")
		return
	}

	if err := a.client.Set(context.Background(), key, value, ttl).Err(); err != nil {
		log.Error().Err(errors.New("ERROR SAVING REFRESH TOKEN ON REDIS : " + err.Error())).Msg("")
		return
	}

	log.Info().Msg("SUCCESS SAVE REFRESH TOKEN WITH KEY " + key)
}

func (a *authRedisRepository) GetRefreshToken(id string) (*models.RefreshTokenRedis, error) {
	key := fmt.Sprintf("%s_refresh-token", id)

	val, err := a.client.Get(context.Background(), key).Result()
	if err != nil {
		log.Error().Err(errors.New("ERROR GETTING REFRESH TOKEN ON REDIS : " + err.Error())).Msg("")
		return nil, nil
	}

	if val == "" {
		log.Error().Err(fmt.Errorf("ERROR REDIS KEY %s : DATA IS EMPTY", key))
		return nil, nil
	}

	var refreshtoken *models.RefreshTokenRedis
	if err := json.Unmarshal([]byte(val), &refreshtoken); err != nil {
		log.Error().Err(fmt.Errorf("ERROR UNMARSHAL REFRESH TOKEN REDIS THAT HAVING KEY %s : %s", key, err.Error())).Msg("")
		return nil, err
	}

	return refreshtoken, nil
}

func (a *authRedisRepository) DeleteRefreshToken(id string) {
	key := fmt.Sprintf("%s_refresh-token", id)

	if err := a.client.Del(context.Background(), key).Err(); err != nil {
		log.Error().Err(errors.New("ERROR DELETE REFRESH TOKEN ON REDIS : " + err.Error())).Msg("")
	}
}

func (a *authRedisRepository) IncLoginFailedCounter(id string) {
	key := fmt.Sprintf("%s_login-failed", id)

	val, err := a.client.Get(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Error().Err(errors.New("ERROR GETTING LOGIN FAILED COUNTER ON REDIS : " + err.Error())).Msg("")
		return
	}

	if val == "" {
		val = ""
	}

	counter, err := strconv.Atoi(val)
	if err != nil {
		log.Error().Err(fmt.Errorf("ERROR REDIS VALUE OF %s ON REDIS IS NOT INTERGER : %v", key, err.Error())).Msg("")
		return
	}

	counter++
	if err := a.client.Set(context.Background(), key, counter, helpers.LoginFailedTTL).Err(); err != nil {
		log.Error().Err(errors.New("ERROR SAVING LOGIN FAILED COUNTER ON REDIS : " + err.Error())).Msg("")
	}
}

func (a *authRedisRepository) IsLoginBlocked(id string) (bool, error) {
	key := fmt.Sprintf("%s_login-failed", id)

	val, err := a.client.Get(context.Background(), key).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		log.Error().Err(errors.New("ERROR GETTING LOGIN FAILED COUNTER ON REDIS : " + err.Error())).Msg("")
		return true, err
	}

	if val == "" {
		val = "0"
	}

	counter, err := strconv.Atoi(val)
	if err != nil {
		log.Error().Err(fmt.Errorf("ERROR REDIS VALUE OF %s ON REDIS IS NOT INTERGER : %v", key, err.Error())).Msg("")
		return true, err
	}

	if counter >= helpers.MaxLoginFailed {
		return true, nil
	}

	return false, nil
}

func (a *authRedisRepository) SaveToken(id string, token string, ttl time.Duration) error {
	key := fmt.Sprintf("%v_token", id)

	if err := a.client.Set(context.Background(), key, token, ttl).Err(); err != nil {
		log.Error().Err(errors.New("ERROR SAVE TOKEN : " + err.Error())).Msg("")
		return err
	}

	log.Info().Msg(fmt.Sprintf("SUCCESS SAVE TOKEN KEY %v", key))
	return nil
}

func (a *authRedisRepository) GetToken(id string) (string, error) {
	key := fmt.Sprintf("%v_token", id)

	val, err := a.client.Get(context.Background(), key).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", nil
		}
		log.Error().Err(errors.New("ERROR GETTING TOKEN REDIS : " + err.Error())).Msg("")
		return "", err
	}

	return val, nil
}

func (a *authRedisRepository) DeleteToken(id string) {
	key := fmt.Sprintf("%v_token", id)

	if err := a.client.Del(context.Background(), key).Err(); err != nil {
		log.Error().Err(errors.New("ERROR DELETE TOKEN REDIS : " + err.Error()))
		return
	}

	log.Info().Msg(fmt.Sprintf("SUCCESS DELETE TOKEN KEY %v", key))
}

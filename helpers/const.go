package helpers

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/phuslu/log"
)

var (
	DefaultTimeFormat           = "2006-01-02 15:04:05"
	MaxLoginFailed              = 3
	LoginFailedTTL              time.Duration
	LoginExpiration             time.Duration
	RefreshTokenExpiration      time.Duration
	RootLocationRedisExpiration time.Duration
)

func LoadConstData() {
	maxLoginFailed, err := strconv.Atoi(os.Getenv("MAX_LOGIN_FAILED"))
	if err != nil {
		log.Error().Err(errors.New("MAX_LOGIN_FAILED IS NOT A NUMBER, CHECK ENV")).Msg("")
		os.Exit(1)
	}

	loginFailedTTL, err := strconv.Atoi(os.Getenv("LOGIN_FAILED_TTL_IN_MINUTES"))
	if err != nil {
		log.Error().Err(errors.New("LOGIN_FAILED_TTL_IN_MINUTES IS NOT A NUMBER, CHECK ENV")).Msg("")
		os.Exit(1)
	}

	loginExpirationTTL, err := strconv.Atoi(os.Getenv("LOGIN_EXPIRATION_IN_MINUTES"))
	if err != nil {
		log.Error().Err(errors.New("LOGIN_EXPIRATION_IN_HOURS IS NOT A NUMBER, CHECK ENV")).Msg("")
		os.Exit(1)
	}

	refreshTokenExpirationTTL, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRATION_IN_MINUTES"))
	if err != nil {
		log.Error().Err(errors.New("REFRESH_TOKEN_EXPIRATION_IN_HOURS IS NOT A NUMBER, CHECK ENV")).Msg("")
		os.Exit(1)
	}

	rootLocationRedisExpirationTTL, err := strconv.Atoi(os.Getenv("ROOT_LOCATION_REDIS_EXPIRATION_IN_HOURS"))
	if err != nil {
		log.Error().Err(errors.New("ROOT_LOCATION_REDIS_EXPIRATION_IN_HOURS IS NOT A NUMBER, CHECK ENV")).Msg("")
		os.Exit(1)
	}

	MaxLoginFailed = maxLoginFailed
	LoginFailedTTL = time.Minute * time.Duration(loginFailedTTL)
	LoginExpiration = time.Minute * time.Duration(loginExpirationTTL)
	RefreshTokenExpiration = time.Minute * time.Duration(refreshTokenExpirationTTL)
	RootLocationRedisExpiration = time.Hour * time.Duration(rootLocationRedisExpirationTTL)
}

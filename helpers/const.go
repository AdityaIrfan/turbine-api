package helpers

import "time"

const (
	DefaultTimeFormat           = "2006-01-02 15:04:05"
	MaxLoginFailed              = 3
	LoginFailedTTL              = time.Duration(time.Minute * 10)
	LoginExpiration             = time.Hour * 2
	RefreshTokenExpiration      = time.Hour * 12
	RootLocationRedisExpiration = time.Hour * 6
)

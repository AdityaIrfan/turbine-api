package helpers

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/phuslu/log"
)

func GenerateToken(userId string, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"Id":  userId,
			"Exp": expiration.Unix(),
		})

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("TOKEN_PRIVATE_KEY")))
	if err != nil {
		log.Error().Err(errors.New("ERROR PARSING PRIVATE KEY TOKEN FROM PEM : " + err.Error())).Msg("")
		return "", err
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Error().Err(errors.New("ERROR SIGNED STRING TOKEN : " + err.Error())).Msg("")
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	verifyKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("TOKEN_PUBLIC_KEY")))

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, _ = token.Method.(*jwt.SigningMethodRSA)

		return verifyKey, nil
	})
}

func GenerateRefreshToken(userId string, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"Id":  userId,
			"Exp": expiration.Unix(),
		})

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(os.Getenv("REFRESH_TOKEN_PRIVATE_KEY")))
	if err != nil {
		log.Error().Err(errors.New("ERROR PARSING PRIVATE KEY REFRESH TOKEN FROM PEM : " + err.Error())).Msg("")
		return "", err
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Error().Err(errors.New("ERROR SIGNED STRING REFRESH TOKEN : " + err.Error())).Msg("")
		return "", err
	}

	return tokenString, nil
}

func VerifyRefreshToken(refreshTokenString string) (*jwt.Token, error) {
	verifyKey, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(os.Getenv("REFRESH_TOKEN_PUBLIC_KEY")))

	return jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		_, _ = token.Method.(*jwt.SigningMethodRSA)

		return verifyKey, nil
	})
}

func IsTokenExpired(token *jwt.Token) bool {
	claims := token.Claims.(jwt.MapClaims)
	expired, ok := claims["Exp"].(float64)
	if !ok {
		log.Error().Err(errors.New("EROR GETTING EXPIRATION VALUE FROM TOKEN : CLAIMS IS NOT EXIST OR VALUE IS NOT int64")).Msg("")
		return true
	}

	return time.Now().After(time.Unix(int64(expired), 0))
}

func IsRefreshTokenActive(token *jwt.Token) bool {
	claims := token.Claims.(jwt.MapClaims)
	expired, ok := claims["Active"].(float64)
	if !ok {
		log.Error().Err(errors.New("EROR GETTING ACTIVE VALUE FROM TOKEN : CLAIMS IS NOT EXIST OR VALUE IS NOT int64")).Msg("")
		return true
	}

	return time.Now().After(time.Unix(int64(expired), 0))
}

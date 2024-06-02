package helpers

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/phuslu/log"
)

var tokenPrivateKey = []byte("-----BEGIN RSA PRIVATE KEY-----\nMIICXQIBAAKBgQDKwLr9Wm6ybMlsscAjckaUBClET4JGdPNew0yWG7qNkHX+QR9y\nhnAZ8Dve53pynZ55BrCMi24JX6prdc3lYsTWhojVuIi8+EzPg6I4puTxkIbJIsvR\n677pWtHOCVCmNuQJ92F6Eb+wdWe3h99A7wL78cBLm4b+fVMYD7DOfFwKKQIDAQAB\nAoGBAJhaZqhGMfRFJr5Uewqnrj0/Ozsne78x1kaY/o21XGNP8DoT+Wt3dfknufVv\nf2Vs504OJtAVXaQRjN+e8A60Py+WD8bMZ8XIYo5pl/OqASAq1m4/MomtsReA2/VL\n8P6ORTey2wSEQUbHFXJ4tWE/CXE8jdguGu7Nv4xJlre4A9rhAkEA5vNdedu9R9NP\nKmG3B+A1MO8MaYK3Zs6oSGsW6clHobYPfq8dvPcoeMhdynnLyk3R8jwoNErdHmMk\nddeOjxpeHQJBAOC+awzknDZM3SOfnAzN7vmI599SxJ6v9Eq5MTPlAr3KZe6np74q\ni/Tdg+Xh1aToJgF5vNmCsvFtGPNhaNrojn0CQEbnABOhOoMKhItmZGKumqXjPdRG\npTeSymcxOV+cw7kJw8gIywBwgKRUHzdCHSaGraXIgi9LrIbfuRnUi5ezaKECQQCj\ngCRlJtO2dUjUJ8PhVNgsVbtKru43/A4fZoczF8JczKhHbVUNdeqH47eXQCqrY/By\nVlxbaUhBd3sVZKJhz5oJAkBQRUigHhCRY7Xu+bojkFKfUV52bRHwAPx6V8pJCsU3\ni9cIBsH9AQC7RgspDo54i1MyIfVwF89sHfljhIlngySd\n-----END RSA PRIVATE KEY-----")
var tokenPublicKey = []byte("-----BEGIN PUBLIC KEY-----\nMIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDKwLr9Wm6ybMlsscAjckaUBClE\nT4JGdPNew0yWG7qNkHX+QR9yhnAZ8Dve53pynZ55BrCMi24JX6prdc3lYsTWhojV\nuIi8+EzPg6I4puTxkIbJIsvR677pWtHOCVCmNuQJ92F6Eb+wdWe3h99A7wL78cBL\nm4b+fVMYD7DOfFwKKQIDAQAB\n-----END PUBLIC KEY-----")

var refreshPrivateKey = []byte("-----BEGIN RSA PRIVATE KEY-----\nMIICWgIBAAKBgF9eZRnM2gFD4aGTq5fjWVTAhmBHrD/mB6KMPrZjvhyyRDX9ZG+k\n1eyay6a7Bwhk19lKZcS25A89aSKAdSObzWXhSAGM+eL9VJlnhqDeXfls2pj1TyUq\nbDDgllK1LaBSGJgLUgs5TUKosJq2q75GalZrclVsd9H4lV/LWFqThq2jAgMBAAEC\ngYBYhCiaVcRSQEr/ySwPAfk77KXMXznE5SBZAoqChvlBcURWSCYcaYYU4Wf8DMNu\nSwr8p0zl1ErAymL2nvwaXNW6gkALPNayFz7Eru22YWqFcNvXTxjmAdy1wGDRNwi+\nhHMkbADW+1WqXKL64EbFGp/NsID93afAitL5EMlPAozKAQJBAK2fN4NQQ8IgRvx1\nWyP5UrtwpXDnJgRR8cUFcu5YoaTzGjqrwlh4U2fz/x0m01X1h5KeLVYd+rpK/FGo\nMH9hE2MCQQCMnj6g9iwcG5NkRey6gdjC0LDKUYSmhTavmK22Qj5FgddPX40OgCqe\nMgeLcFLRdd4DcSkesagHaaIzfuhs4bDBAkA2o+HjmJIKeP/+GazaMG/h/3yBgK1N\nNMDCwYk/C3Orpro9dqqODygokfhao0plRgUpllAsRvkOQeUQib7hh5qDAkA7x4ZO\nfXjxFhQJ2+QwwcS5xWhzCka/WACQk/K9ednpSLKU7sUTtg7oI9KrR7wdieMxSWk2\nwEXzqMeo5rm+mA/BAkB6Tx6+BN40pz/KW1PJciGPjAR/xG66W9dYEA2U0HifGZAa\nrO0X5XyQUuwFBq7KUaSaA3n/A7gZxVOwkNhcNoi7\n-----END RSA PRIVATE KEY-----")
var refreshPublicKey = []byte("-----BEGIN PUBLIC KEY-----\nMIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgF9eZRnM2gFD4aGTq5fjWVTAhmBH\nrD/mB6KMPrZjvhyyRDX9ZG+k1eyay6a7Bwhk19lKZcS25A89aSKAdSObzWXhSAGM\n+eL9VJlnhqDeXfls2pj1TyUqbDDgllK1LaBSGJgLUgs5TUKosJq2q75GalZrclVs\nd9H4lV/LWFqThq2jAgMBAAE=\n-----END PUBLIC KEY-----")

func GenerateToken(userId string, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"UserId": userId,
			"Exp":    expiration.Unix(),
		})

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(tokenPrivateKey)
	if err != nil {
		log.Error().Err(errors.New("ERROR PARSING PRIVATE KEY TOKEN FROM PEM : " + err.Error()))
		return "", err
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Error().Err(errors.New("ERROR SIGNED STRING TOKEN : " + err.Error()))
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (*jwt.Token, error) {
	verifyKey, _ := jwt.ParseRSAPublicKeyFromPEM(tokenPublicKey)

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, _ = token.Method.(*jwt.SigningMethodRSA)

		return verifyKey, nil
	})
}

func GenerateRefreshToken(userId string, expiration time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256,
		jwt.MapClaims{
			"UserId": userId,
			"Exp":    expiration.Unix(),
		})

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(refreshPrivateKey)
	if err != nil {
		log.Error().Err(errors.New("ERROR PARSING PRIVATE KEY REFRESH TOKEN FROM PEM : " + err.Error()))
		return "", err
	}

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		log.Error().Err(errors.New("ERROR SIGNED STRING REFRESH TOKEN : " + err.Error()))
		return "", err
	}

	return tokenString, nil
}

func VerifyRefreshToken(refreshTokenString string) (*jwt.Token, error) {
	verifyKey, _ := jwt.ParseRSAPublicKeyFromPEM(refreshPublicKey)

	return jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		_, _ = token.Method.(*jwt.SigningMethodRSA)

		return verifyKey, nil
	})
}

func IsTokenExpired(token *jwt.Token) bool {
	claims := token.Claims.(jwt.MapClaims)
	expired, ok := claims["Exp"].(float64)
	if !ok {
		log.Error().Err(errors.New("EROR GETTING EXPIRATION VALUE FROM TOKEN : CLAIMS IS NOT EXIST OR VALUE IS NOT int64"))
		return true
	}

	if time.Now().After(time.Unix(int64(expired), 0)) {
		return true
	}

	return false
}

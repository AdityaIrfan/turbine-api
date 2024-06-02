package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	mathRand "math/rand"
	"os"

	"github.com/phuslu/log"
	"golang.org/x/crypto/bcrypt"
)

func RandomInt(min int, max int) int {
	return min + mathRand.Intn(max-min)
}

func RandomByte(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)

	return base64.URLEncoding.EncodeToString(b)
}

func Encrypt(plainText string) string {
	key, _ := hex.DecodeString(os.Getenv("APP_KEY"))
	block, _ := aes.NewCipher(key)
	aesGCM, _ := cipher.NewGCM(block)
	nonce := make([]byte, aesGCM.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	return fmt.Sprintf("%x", aesGCM.Seal(nonce, nonce, []byte(plainText), nil))
}

// return salt, hash, error
func GenerateHashAndSalt(key string) (string, string, error) {
	salt := RandomByte(RandomInt(1024, 2048))

	hash, err := bcrypt.GenerateFromPassword([]byte(key+salt), 7)
	if err != nil {
		log.Error().Err(errors.New("ERROR PASSWORD GENERATE HASH AND SALT : " + err.Error()))
		return "", "", err
	}

	return Encrypt(salt), Encrypt(string(hash)), nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890@#&!")

func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[mathRand.Intn(len(letters))]
	}
	return string(b)
}

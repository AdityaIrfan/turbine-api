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
	key, _ := hex.DecodeString("751dc352fcd484748c31b7d8562cf9778843492ebffa26de11e59883865bba82")
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println(err)
	}
	aesGCM, _ := cipher.NewGCM(block)
	nonce := make([]byte, aesGCM.NonceSize())
	_, _ = io.ReadFull(rand.Reader, nonce)

	return fmt.Sprintf("%x", aesGCM.Seal(nonce, nonce, []byte(plainText), nil))
}

// return salt, hash, error
func GenerateHashAndSalt(key string) (string, string, error) {
	salt := RandomByte(RandomInt(32, 64))

	hash, err := bcrypt.GenerateFromPassword([]byte(key+salt), 7)
	if err != nil {
		log.Error().Err(errors.New("ERROR PASSWORD GENERATE HASH AND SALT : " + err.Error())).Msg("")
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

// Decrypt is function to decrypt chiper text
func Decrypt(encryptedString string) (string, error) {
	key, err := hex.DecodeString("751dc352fcd484748c31b7d8562cf9778843492ebffa26de11e59883865bba82")
	if err != nil {
		return "", err
	}

	enc, err := hex.DecodeString(encryptedString)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, cipherText := enc[:nonceSize], enc[nonceSize:]

	plainText, err := aesGCM.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

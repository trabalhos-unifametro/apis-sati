package utils

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
)

var (
	MANDRILL string
	FROM     string
)

func IsEmpty(value string) bool {
	if len(strings.TrimSpace(value)) == 0 {
		return true
	}
	return false
}

func EncryptPassword(password string) string {
	passwordByte := []byte(password)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordByte, bcrypt.DefaultCost)
	if err != nil {
		LogMessage{Title: "[MAIN>UTILS] Error on EncryptPassword()", Body: err.Error()}.Error()
	}
	return string(hashedPassword)
}

func ComparePassword(password1, password2 string) error {
	return bcrypt.CompareHashAndPassword(
		[]byte(password1),
		[]byte(password2),
	)
}

func GenerateCode(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func LoadENVs() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		LogMessage{Title: "[MAIN>UTILS] Error on LoadENVs()", Body: err.Error()}.Error()
	}
	MANDRILL = viper.GetString("MANDRILL")
	FROM = viper.GetString("FROM")
}

func CheckToSend(emails string) string {
	env := viper.GetString("ENV")
	if env == "DEVELOPMENT" {
		return strings.ToLower(viper.GetString("EMAIL_DEV"))
	}
	return strings.ToLower(emails)
}

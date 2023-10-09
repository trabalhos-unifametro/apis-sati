package utils

import (
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strconv"
	"strings"
)

var (
	SMTP     string
	USERNAME string
	PASSWORD string
	FROM     string
	PORT     int
)

func IsEmpty(value string) bool {
	if len(strings.TrimSpace(value)) == 0 {
		return true
	}
	return false
}

func ToComparePassword(password1, password2 string) error {
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
		LogMessage{Title: "[MAIN] Error on *User.FindUserEmail()", Body: err.Error()}.Error()
	}
	SMTP = viper.GetString("SMTP")
	USERNAME = viper.GetString("USERNAME")
	PASSWORD = viper.GetString("PASSWORD_EMAIL")
	FROM = viper.GetString("FROM")
	PORT, _ = strconv.Atoi(viper.GetString("PORT"))
}

func CheckToSend(emails string) string {
	env := viper.GetString("ENV")
	if env == "DEVELOPMENT" {
		return strings.ToLower(viper.GetString("EMAIL_DEV"))
	}
	return strings.ToLower(emails)
}

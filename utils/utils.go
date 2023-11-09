package utils

import (
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
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
	MANDRILL = GodotEnv("MANDRILL")
	FROM = GodotEnv("FROM")
}

func CheckToSend(emails string) string {
	env := GodotEnv("ENV")
	if env == "DEVELOPMENT" {
		return strings.ToLower(GodotEnv("EMAIL_DEV"))
	}
	return strings.ToLower(emails)
}

func GodotEnv(key string) string {
	var err error
	if os.Getenv("env") == "development" {
		return os.Getenv(key)
	} else {
		err = godotenv.Load(".env.production")
		if err != nil {
			LogMessage{Title: "[MAIN>UTILS] Error on LoadENVs()", Body: err.Error()}.Error()
		}
		return os.Getenv(key)
	}
}

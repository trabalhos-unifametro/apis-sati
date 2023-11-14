package utils

import (
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	SMTP     string
	FROM     string
	PORT     int
	USERNAME string
	PASSWORD string
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

func GenerateCode(length int) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())

	code := make([]byte, length)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}

	return string(code)
}

func LoadENVs() {
	SMTP = GodotEnv("SMTP_EMAIL")
	FROM = GodotEnv("FROM_EMAIL")
	USERNAME = GodotEnv("USERNAME_EMAIL")
	PASSWORD = GodotEnv("PASSWORD_EMAIL")
	PORT, _ = strconv.Atoi(GodotEnv("PORT_EMAIL"))
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

package controllers

import (
	"apis-sati/auth"
	"apis-sati/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func ValidateTokenSession(c *fiber.Ctx) (auth.JwtClaim, error) {
	tokenString := c.Get("Authorization")
	jwtWrapper := auth.JwtWrapper{SecretKey: utils.GodotEnv("JWT_SECRET_KEY")}
	claims, err := jwtWrapper.ValidateToken(tokenString)
	if tokenString != "" {
		if err == nil {
			return *claims, err
		} else {
			return auth.JwtClaim{}, errors.New("Você não tem acesso a essa rota!")
		}
	} else {
		return auth.JwtClaim{}, errors.New("Você não tem acesso a essa rota!")
	}
}

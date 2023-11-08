package controllers

import (
	"apis-sati/models"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func ListStatusPatients(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err == nil {
		var status models.Status
		var list []models.ResponseStatus
		t := c.Query("t")
		var table string

		if strings.ToLower(t) != "mr" && strings.ToLower(t) != "u" {
			return c.Status(fiber.StatusBadRequest).JSON("Não foi possível carregar os status.")
		}

		if strings.ToLower(t) == "mr" {
			table = "medical_records"
		} else {
			table = "units"
		}

		if err, list = status.GetListStatus(table); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar os status.")
		}

		return c.Status(fiber.StatusOK).JSON(list)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

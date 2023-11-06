package controllers

import (
	"apis-sati/models"
	"github.com/gofiber/fiber/v2"
)

func GraphicMonthlyDashboard(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err != nil {
		var medicalRecord models.MedicalRecord
		var graphic []models.MonthlyChart

		if err, graphic = medicalRecord.GraphicDashboard(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar o gráfico mensal de ocupação.")
		}

		return c.Status(fiber.StatusOK).JSON(graphic)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

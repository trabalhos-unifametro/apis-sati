package controllers

import (
	"apis-sati/models"
	"github.com/gofiber/fiber/v2"
)

func GraphicMonthlyDashboard(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err == nil {
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

func GetMedicalRecord(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err == nil {
		id, err := c.ParamsInt("id")

		var medicalRecord = models.MedicalRecord{ID: uint(id)}
		var response models.ResponseMedicalRecord

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Informe o ID correto do prontuário.")
		}

		if err, response = medicalRecord.GetMedicalRecordByID(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar os dados do prontuário.")
		}

		return c.Status(fiber.StatusOK).JSON(response)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

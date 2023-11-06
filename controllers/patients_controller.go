package controllers

import (
	"apis-sati/models"
	"github.com/gofiber/fiber/v2"
)

func TotalizatorsPatientsDashboard(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err != nil {
		var unit models.Patient
		var totalizators models.PatientTotalizatorsDashboard

		if err, totalizators = unit.TotalizatorsDashboard(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar os totalizadores dos pacientes.")
		}

		return c.Status(fiber.StatusOK).JSON(totalizators)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

func GraphicPatientsDashboard(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err != nil {
		var unit models.Patient
		var graphic models.PatientGraphicDashboard

		if err, graphic = unit.GraphicDashboard(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar o gráfico de pacientes.")
		}

		return c.Status(fiber.StatusOK).JSON(graphic)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

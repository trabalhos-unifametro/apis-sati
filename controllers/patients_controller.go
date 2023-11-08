package controllers

import (
	"apis-sati/models"
	"github.com/gofiber/fiber/v2"
)

func TotalizatorsPatientsDashboard(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err == nil {
		var patient models.Patient
		var totalizators models.PatientTotalizatorsDashboard

		if err, totalizators = patient.TotalizatorsDashboard(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar os totalizadores dos pacientes.")
		}

		return c.Status(fiber.StatusOK).JSON(totalizators)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

func GraphicPatientsDashboard(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err == nil {
		var patient models.Patient
		var graphic models.PatientGraphicDashboard

		if err, graphic = patient.GraphicDashboard(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar o gráfico de pacientes.")
		}

		return c.Status(fiber.StatusOK).JSON(graphic)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

func TotalizatorsPatients(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err == nil {
		var patient models.Patient
		var totalizators models.PatientTotalizators
		patientName := c.Query("patient_name")
		situation := c.Query("situation_patient")
		statusID := c.QueryInt("status_id")

		if err, totalizators = patient.Totalizators(patientName, situation, statusID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar os totalizadores dos pacientes.")
		}

		return c.Status(fiber.StatusOK).JSON(totalizators)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

func ListPatients(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err == nil {
		patientName := c.Query("patient_name")
		situationPatient := c.Query("situation_patient")
		sortByPatient := c.Query("sort_by_patient")
		statusID := c.QueryInt("status_id")

		var patient = models.Patient{}
		var list []models.ResponsePatient

		if err, list = patient.GetListPatients(patientName, situationPatient, sortByPatient, statusID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar os pacientes.")
		}

		return c.Status(fiber.StatusOK).JSON(list)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

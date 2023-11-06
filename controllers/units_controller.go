package controllers

import (
	"apis-sati/models"
	"github.com/gofiber/fiber/v2"
)

func TotalizatorsUnitsDashboard(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err != nil {
		var unit models.Unit
		var totalizators models.UnitTotalizatorsDashboard
		unitName := c.Query("unit_name")
		occupation := c.Query("occupation")

		if err, totalizators = unit.TotalizatorsDashboard(unitName, occupation); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar os totalizadores das unidades.")
		}

		return c.Status(fiber.StatusOK).JSON(totalizators)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

func GraphicUnitsDashboard(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err != nil {
		var unit models.Unit
		var graphic models.UnitGraphicDashboard

		if err, graphic = unit.GraphicDashboard(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar o gráfico de unidades.")
		}

		return c.Status(fiber.StatusOK).JSON(graphic)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

func ListUnits(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err != nil {
		var unit models.Unit
		var list []models.ResponseUnit
		unitName := c.Query("unit_name")
		occupation := c.Query("occupation")
		sortByUnit := c.Query("sort_by_unit")
		sortByOccupation := c.Query("sort_by_occupation")

		if err, list = unit.GetListUnits(unitName, occupation, sortByUnit, sortByOccupation); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar a lista de unidades.")
		}

		return c.Status(fiber.StatusOK).JSON(list)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

func TotalizatorsUnit(c *fiber.Ctx) error {
	if _, err := ValidateTokenSession(c); err != nil {
		id, err := c.ParamsInt("id")
		patientName := c.Query("patient_name")
		situationPatient := c.Query("situation_patient")
		var unit = models.Unit{ID: uint(id)}
		var totalizators models.UnitTotalizators

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Informe o ID correto da unidade.")
		}

		if err, totalizators = unit.Totalizators(patientName, situationPatient); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar os totalizadores da unidade.")
		}

		return c.Status(fiber.StatusOK).JSON(totalizators)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

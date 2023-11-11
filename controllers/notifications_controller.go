package controllers

import (
	"apis-sati/models"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

func ListNotifications(c *fiber.Ctx) error {
	if session, err := ValidateTokenSession(c); err == nil {
		userID, _ := c.ParamsInt("userID")
		var notification models.Notification
		var list []models.ResponseNotification

		if session.UserID != userID {
			return c.Status(fiber.StatusBadRequest).JSON("Usuário informado é desconhecido.")
		}

		if err, list = notification.GetListNotificationsByUserID(userID); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao carregar as notificações.")
		}

		return c.Status(fiber.StatusOK).JSON(list)
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

func ReadNotifications(c *fiber.Ctx) error {
	if session, err := ValidateTokenSession(c); err == nil {
		userID, _ := c.ParamsInt("userID")
		var notification models.Notification
		body := c.Body()

		if session.UserID != userID {
			return c.Status(fiber.StatusBadRequest).JSON("Usuário informado é desconhecido.")
		}

		if err := json.Unmarshal(body, &notification); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Estrutura de dados incorreta!")
		}

		if success := notification.ReadNotificationByUserID(userID); !success {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao salvar notificações como lidas.")
		}

		return c.Status(fiber.StatusOK).JSON("Notificação foi marcada como lida com sucesso!")
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

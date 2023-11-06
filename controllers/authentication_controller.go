package controllers

import (
	"apis-sati/auth"
	"apis-sati/emails"
	"apis-sati/models"
	"apis-sati/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func SignIn(c *fiber.Ctx) error {
	body := c.Body()
	user := models.User{}

	if err := json.Unmarshal(body, &user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Estrutura de dados incorreta!")
	}

	if utils.IsEmpty(user.Email) || utils.IsEmpty(user.Password) {
		return c.Status(fiber.StatusBadRequest).JSON("É obrigatório informar email e senha.")
	}

	if err := user.FindByEmail(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Usuário não cadastrado na base.")
	}

	errPassword := utils.ComparePassword(user.PasswordDigest, user.Password)

	if user.ID > 0 {
		if errPassword == nil {
			jwtWrapper := auth.JwtWrapper{
				SecretKey:       viper.GetString("JWT_SECRET_KEY"),
				Issuer:          "AuthService",
				ExpirationHours: 120,
			}

			generatedToken, _ := jwtWrapper.GenerateToken(int(user.ID))
			response := models.DataUser{
				ID:        user.ID,
				Name:      user.Name,
				Email:     user.Email,
				Cellphone: user.Cellphone,
				Role:      user.Role,
				CreatedAt: user.CreatedAt,
				Token:     generatedToken,
			}
			return c.Status(fiber.StatusOK).JSON(response)
		} else {
			return c.Status(fiber.StatusBadRequest).JSON("Senha incorreta.")
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON("Usuário não cadastrado na base.")
	}
}

func GenerateCodeRecoverPassword(c *fiber.Ctx) error {
	body := c.Body()
	user := models.User{}

	if err := json.Unmarshal(body, &user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Estrutura de dados incorreta!")
	}

	if utils.IsEmpty(user.Email) {
		return c.Status(fiber.StatusBadRequest).JSON("É obrigatório informar o email.")
	}

	if err := user.FindByEmail(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Usuário não cadastrado na base.")
	}

	user.CodeRecovery = utils.GenerateCode(7)

	if user.SaveCodeRecover() {
		if err := emails.SendEmailCodeRandom(user); err == nil {
			return c.Status(fiber.StatusOK).JSON("O código foi enviado ao seu email!")
		} else {
			return c.Status(fiber.StatusBadRequest).JSON("Erro ao tentar enviar código.")
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON("Erro ao tentar gerar código.")
	}
}

func ConfirmCode(c *fiber.Ctx) error {
	body := c.Body()
	user := models.User{}

	if err := json.Unmarshal(body, &user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Estrutura de dados incorreta!")
	}

	if utils.IsEmpty(user.Email) || utils.IsEmpty(user.CodeRecovery) {
		return c.Status(fiber.StatusBadRequest).JSON("É obrigatório informar o email e o código.")
	}

	if err := user.FindByEmail(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Usuário não cadastrado na base.")
	}

	if user.ConfirmCodeRecover() {
		return c.Status(fiber.StatusOK).JSON("Código validado com sucesso!")
	} else {
		return c.Status(fiber.StatusBadRequest).JSON("Erro ao tentar validar código.")
	}
}

func RecoverPassword(c *fiber.Ctx) error {
	body := c.Body()
	dataUser := models.DataUser{}
	user := models.User{}

	if err := json.Unmarshal(body, &dataUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Estrutura de dados incorreta!")
	}
	user.Email = dataUser.Email
	user.CodeRecovery = dataUser.CodeRecovery

	if isValid, message := dataUser.ValidationRecoverPassword(); !isValid {
		return c.Status(fiber.StatusBadRequest).JSON(message)
	}

	if err := user.FindByEmailAndCodeRecovery(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON("Usuário não cadastrado na base.")
	}

	user.PasswordDigest = utils.EncryptPassword(dataUser.NewPassword)
	if user.UpdatePassword() {
		if err := emails.SuccessfulRecoverPassword(user); err == nil {
			return c.Status(fiber.StatusOK).JSON("Nova senha atualizada com sucesso!")
		} else {
			return c.Status(fiber.StatusBadRequest).JSON("Erro ao tentar enviar email de nova senha.")
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON("Erro ao tentar atualizar nova senha.")
	}
}

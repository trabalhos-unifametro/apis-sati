package controllers

import (
	"apis-sati/auth"
	"apis-sati/emails"
	"apis-sati/models"
	"apis-sati/utils"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
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
				SecretKey:       utils.GodotEnv("JWT_SECRET_KEY"),
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
	if user.ResetPassword() {
		if err := emails.SuccessfulRecoverPassword(user); err == nil {
			return c.Status(fiber.StatusOK).JSON("Nova senha atualizada com sucesso!")
		} else {
			return c.Status(fiber.StatusBadRequest).JSON("Erro ao tentar enviar email de nova senha.")
		}
	} else {
		return c.Status(fiber.StatusBadRequest).JSON("Erro ao tentar atualizar nova senha.")
	}
}

func UpdateDataUser(c *fiber.Ctx) error {
	if session, err := ValidateTokenSession(c); err == nil {
		body := c.Body()
		user := models.User{}

		if err := json.Unmarshal(body, &user); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Estrutura de dados incorreta!")
		}

		if uint(session.UserID) != user.ID {
			return c.Status(fiber.StatusBadRequest).JSON("Identificador do usuário está incorreto!")
		}

		if exists := user.FindUserByEmailAndNotID(); exists {
			return c.Status(fiber.StatusBadRequest).JSON("Email já está em uso, escolha outro.")
		}

		if utils.IsEmpty(user.Email) || utils.IsEmpty(user.Cellphone) {
			return c.Status(fiber.StatusBadRequest).JSON("É obrigatório informar email e telefone.")
		}

		if err = user.Update(); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Ocorreu um erro ao atualizar os dados do usuário.")
		}

		return c.Status(fiber.StatusOK).JSON("Dados atualizados com sucesso!")
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

func ChangePassword(c *fiber.Ctx) error {
	if session, err := ValidateTokenSession(c); err == nil {
		body := c.Body()
		dataUser := models.DataUser{}
		user := models.User{}

		if err := json.Unmarshal(body, &dataUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON("Estrutura de dados incorreta!")
		}

		if uint(session.UserID) != dataUser.ID {
			return c.Status(fiber.StatusBadRequest).JSON("Usuário ID incorreto!")
		}

		user.ID = dataUser.ID

		if exists := user.FindUserByID(); exists == false {
			return c.Status(fiber.StatusBadRequest).JSON("Usuário não foi encontrado.")
		}

		if isValid, message := dataUser.ValidationChangePassword(); !isValid {
			return c.Status(fiber.StatusBadRequest).JSON(message)
		}

		user.PasswordDigest = utils.EncryptPassword(dataUser.NewPassword)

		if success := user.UpdatePassword(); success {
			return c.Status(fiber.StatusOK).JSON("Nova senha salva com sucesso!")
		} else {
			return c.Status(fiber.StatusBadRequest).JSON("Erro ao tentar atualizar senha.")
		}
	} else {
		return c.Status(fiber.StatusUnauthorized).JSON("Você não tem acesso a essa rota ")
	}
}

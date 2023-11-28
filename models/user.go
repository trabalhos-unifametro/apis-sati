package models

import (
	"apis-sati/database"
	"apis-sati/utils"
	"gorm.io/gorm"
	"regexp"
	"time"
)

type User struct {
	gorm.Model
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `gorm:"type:varchar(200)" rql:"filter" json:"name"`
	Role           string         `gorm:"type:varchar(100)" rql:"filter" json:"role"`
	Email          string         `gorm:"type:varchar(100)" rql:"filter" json:"email"`
	Cellphone      string         `gorm:"type:varchar(16)" rql:"filter" json:"cellphone"`
	CodeRecovery   string         `gorm:"type:varchar(7)" rql:"filter" json:"code_recovery"`
	Password       string         `gorm:"type:varchar(100)" rql:"filter" json:"password"`
	PasswordDigest string         `gorm:"type:varchar(100)" rql:"filter" json:"password_digest"`
	ExpirationCode time.Time      `json:"expiration_code"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type DataUser struct {
	ID                 uint      `json:"id,omitempty"`
	Name               string    `json:"name,omitempty"`
	Role               string    `json:"role,omitempty"`
	Email              string    `json:"email,omitempty"`
	Cellphone          string    `json:"cellphone,omitempty"`
	CurrentPassword    string    `json:"current_password,omitempty"`
	NewPassword        string    `json:"new_password,omitempty"`
	ConfirmNewPassword string    `json:"confirm_new_password,omitempty"`
	CodeRecovery       string    `json:"code_recovery,omitempty"`
	Token              string    `json:"token,omitempty"`
	CreatedAt          time.Time `json:"created_at,omitempty"`
}

func (u *User) FindByEmail() error {
	db := database.OpenConnection()
	err := db.Table("sati.users").
		Select("id, name, role, email, cellphone, password_digest, created_at").
		Where("email ILIKE ?", u.Email).Find(&u).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on *User.FindUserEmail()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on database.CloseConnection(db) > *User.FindUserEmail()", Body: err.Error()}.Error()
	}
	return err
}

func (u *User) FindByEmailAndCodeRecovery() error {
	db := database.OpenConnection()
	err := db.Table("sati.users").
		Select("id, name, role, email, cellphone, password_digest, password, created_at").
		Where("email ILIKE ? AND code_recovery = ? AND to_char(expiration_code, 'YYYY-MM-DD HH24:MI') >= to_char(current_timestamp, 'YYYY-MM-DD HH24:MI')", u.Email, u.CodeRecovery).
		Find(&u).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on *User.FindByEmailAndCodeRecovery()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on database.CloseConnection(db) > *User.FindByEmailAndCodeRecovery()", Body: err.Error()}.Error()
	}
	return err
}

func (u *User) SaveCodeRecover() bool {
	db := database.OpenConnection()
	var codeHasSave bool
	err := db.Table("sati.users").
		Where("id = ?", u.ID).
		Updates(map[string]interface{}{"code_recovery": u.CodeRecovery, "expiration_code": time.Now().Add(time.Minute * 5)}).Error
	if err != nil {
		codeHasSave = false
	} else {
		codeHasSave = true
	}
	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on database.CloseConnection(db) > *User.SaveCodeRecover()", Body: err.Error()}.Error()
	}

	return codeHasSave
}

func (u *User) ConfirmCodeRecover() bool {
	db := database.OpenConnection()
	err := db.Table("sati.users").
		Select("id").
		Where("email ILIKE ? AND code_recovery = ? AND to_char(expiration_code, 'YYYY-MM-DD HH24:MI') >= to_char(current_timestamp, 'YYYY-MM-DD HH24:MI')", u.Email, u.CodeRecovery).
		Find(&u).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on *User.ConfirmCodeRecover()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on database.CloseConnection(db) > *User.ConfirmCodeRecover()", Body: err.Error()}.Error()
	}
	if u.ID == 0 {
		return false
	}

	return true
}

func (u DataUser) ValidationRecoverPassword() (bool, string) {

	regLowerCase := regexp.MustCompile(`[A-Z]`)
	regUpperCase := regexp.MustCompile(`[A-Z]`)
	regSpecialCharacters := regexp.MustCompile(`[!@#$%^&*()_+\-=[\]{};':"|,.<>/?~]`)
	regNumbers := regexp.MustCompile(`[0-9]`)

	if utils.IsEmpty(u.Email) {
		return false, "Por favor, informe o email."
	}
	if utils.IsEmpty(u.CodeRecovery) {
		return false, "Por favor, informe o código de verificação."
	}
	if utils.IsEmpty(u.NewPassword) {
		return false, "Por favor, informe a nova senha."
	}
	if utils.IsEmpty(u.ConfirmNewPassword) {
		return false, "Por favor, confirme a nova senha."
	}
	if len(u.NewPassword) < 8 {
		return false, "A nova senha precisa conter no mínimo 8 caracteres."
	}
	if !regLowerCase.MatchString(u.NewPassword) {
		return false, "A nova senha precisa conter no mínimo 1 caractere minúsculo."
	}
	if !regUpperCase.MatchString(u.NewPassword) {
		return false, "A nova senha precisa conter no mínimo 1 caractere maiúsculo."
	}
	if !regSpecialCharacters.MatchString(u.NewPassword) {
		return false, "A nova senha precisa conter no mínimo 1 caractere especial."
	}
	if !regNumbers.MatchString(u.NewPassword) {
		return false, "A nova senha precisa conter no mínimo 1 caractere numérico."
	}
	if u.ConfirmNewPassword != u.NewPassword {
		return false, "As senhas estão diferentes."
	}
	return true, ""
}

func (u *User) ResetPassword() bool {
	db := database.OpenConnection()
	result := db.Table("sati.users").
		Where("id = ? AND code_recovery = ? AND to_char(expiration_code, 'YYYY-MM-DD HH24:MI') >= to_char(current_timestamp, 'YYYY-MM-DD HH24:MI')", u.ID, u.CodeRecovery).
		Updates(map[string]interface{}{"password_digest": u.PasswordDigest, "password": u.PasswordDigest, "updated_at": time.Now()})
	if result.Error != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on *User.ResetPassword()", Body: result.Error.Error()}.Error()
	}
	if err := database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on database.CloseConnection(db) > *User.ResetPassword()", Body: err.Error()}.Error()
	}

	return result.RowsAffected > 0
}

func (u *User) FindUserByEmailAndNotID() bool {
	db := database.OpenConnection()
	var exists bool
	err := db.Table("sati.users").
		Select("(CASE WHEN COUNT(*) > 0 THEN true ELSE false END) as exists").
		Where("email ILIKE ? AND id <> ?", u.Email, u.ID).
		Limit(1).
		Find(&exists).Error

	if err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on *User.FindUserByEmailAndID()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on database.CloseConnection(db) > *User.FindUserByEmailAndID()", Body: err.Error()}.Error()
	}

	return exists
}

func (u *User) Update() error {
	db := database.OpenConnection()

	err := db.Table("sati.users").
		Where("id = ? ", u.ID).
		Updates(map[string]interface{}{"email": u.Email, "cellphone": u.Cellphone, "updated_at": time.Now()}).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on *User.Update()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on database.CloseConnection(db) > *User.Update()", Body: err.Error()}.Error()
	}
	return err
}

func (u *User) UpdatePassword() bool {
	db := database.OpenConnection()
	var success bool
	result := db.Table("sati.users").
		Where("id = ? ", u.ID).
		Updates(map[string]interface{}{"password_digest": u.PasswordDigest, "password": u.PasswordDigest})
	if result.Error != nil || result.RowsAffected == 0 {
		success = false
	} else {
		success = true
	}
	if err := database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on database.CloseConnection(db) > *User.UpdatePassword()", Body: err.Error()}.Error()
	}

	return success
}

func (u *User) FindUserByID() bool {
	db := database.OpenConnection()
	var exists bool
	err := db.Table("sati.users").
		Select("(CASE WHEN COUNT(*) > 0 THEN true ELSE false END) as exists").
		Where("id = ?", u.ID).
		Limit(1).
		Find(&exists).Error

	if err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on *User.FindUserByID()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on database.CloseConnection(db) > *User.FindUserByID()", Body: err.Error()}.Error()
	}

	return exists
}

func (u *User) GetPasswordByID() error {
	db := database.OpenConnection()
	err := db.Table("sati.users").
		Select("password_digest").
		Where("id = ?", u.ID).Find(&u).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on *User.GetPasswordByID()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on database.CloseConnection(db) > *User.GetPasswordByID()", Body: err.Error()}.Error()
	}
	return err
}

func (u DataUser) ValidationChangePassword() (bool, string) {
	user := User{ID: u.ID}
	if err := user.GetPasswordByID(); err != nil {
		return false, "Usuário não encontrado."
	}

	regLowerCase := regexp.MustCompile(`[A-Z]`)
	regUpperCase := regexp.MustCompile(`[A-Z]`)
	regSpecialCharacters := regexp.MustCompile(`[!@#$%^&*()_+\-=[\]{};':"|,.<>/?~]`)
	regNumbers := regexp.MustCompile(`[0-9]`)

	if errPassword := utils.ComparePassword(user.PasswordDigest, u.CurrentPassword); errPassword != nil {
		return false, "Senha atual está incorreta."
	}
	if utils.IsEmpty(u.NewPassword) {
		return false, "Por favor, informe a nova senha."
	}
	if utils.IsEmpty(u.ConfirmNewPassword) {
		return false, "Por favor, confirme a nova senha."
	}
	if len(u.NewPassword) < 8 {
		return false, "A nova senha precisa conter no mínimo 8 caracteres."
	}
	if !regLowerCase.MatchString(u.NewPassword) {
		return false, "A nova senha precisa conter no mínimo 1 caractere minúsculo."
	}
	if !regUpperCase.MatchString(u.NewPassword) {
		return false, "A nova senha precisa conter no mínimo 1 caractere maiúsculo."
	}
	if !regSpecialCharacters.MatchString(u.NewPassword) {
		return false, "A nova senha precisa conter no mínimo 1 caractere especial."
	}
	if !regNumbers.MatchString(u.NewPassword) {
		return false, "A nova senha precisa conter no mínimo 1 caractere numérico."
	}
	if u.ConfirmNewPassword != u.NewPassword {
		return false, "As senhas estão diferentes."
	}
	return true, ""
}

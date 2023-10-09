package models

import (
	"apis-sati/database"
	"apis-sati/utils"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID             uint           `json:"id" gorm:"primaryKey"`
	Name           string         `gorm:"type:varchar(200)" rql:"filter" json:"name"`
	Role           string         `gorm:"type:varchar(100)" rql:"filter" json:"role"`
	Email          string         `gorm:"type:varchar(100)" rql:"filter" json:"email"`
	Cellphone      string         `gorm:"type:varchar(15)" rql:"filter" json:"cellphone"`
	CodeRecovery   string         `gorm:"type:varchar(7)" rql:"filter" json:"code_recovery"`
	Password       string         `gorm:"type:varchar(100)" rql:"filter" json:"password"`
	PasswordDigest string         `gorm:"type:varchar(100)" rql:"filter" json:"password_digest"`
	ExpirationCode time.Time      `json:"expiration_code"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type ResponseUser struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	Email     string    `json:"email"`
	Cellphone string    `json:"cellphone"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) FindByEmail() error {
	db := database.OpenConnection()
	err := db.Table("users").
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

func (u *User) SaveCodeRecover() bool {
	db := database.OpenConnection()
	var codeHasSave bool
	err := db.Table("users").
		Where("id = ?", u.ID).
		Updates(map[string]interface{}{"code_recovery": u.CodeRecovery, "expiration_code": time.Now()}).Error
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
	var dateCurrent = time.Now().Format("2006-01-02")
	var minutes = time.Now().Minute()
	where := fmt.Sprint(`email ILIKE `, u.Email, `
			AND code_recovery = '`, u.CodeRecovery, `'
			AND expiration_code 
			BETWEEN '`, dateCurrent, ' ', minutes-5, `:00' AND '`, dateCurrent, ' ', minutes, `:00'`)
	err := db.Table("users").
		Select("id").
		Where(where).Find(&u).Error
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

func (u *User) UpdateExpirationCode() bool {
	db := database.OpenConnection()
	var success bool
	err := db.Table("users").
		Where("id = ?", u.ID).
		Updates(map[string]interface{}{"expiration_code": time.Now().Add(time.Minute * 30)}).Error
	if err != nil {
		success = false
	} else {
		success = true
	}
	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>USER] Error on database.CloseConnection(db) > *User.UpdateExpirationCode()", Body: err.Error()}.Error()
	}

	return success
}

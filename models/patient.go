package models

import (
	"gorm.io/gorm"
	"time"
)

type Patient struct {
	gorm.Model
	ID          uint           `json:"id" gorm:"primaryKey"`
	MotherName  string         `gorm:"type:varchar(200)" rql:"filter" json:"mother_name"`
	Name        string         `gorm:"type:varchar(200)" rql:"filter" json:"name"`
	Cpf         string         `gorm:"type:varchar(14)" rql:"filter" json:"cpf"`
	Cellphone   string         `gorm:"type:varchar(15)" rql:"filter" json:"cellphone"`
	Telephone   string         `gorm:"type:varchar(14)" rql:"filter" json:"telephone"`
	DateOfBirth time.Time      `rql:"filter" json:"date_of_birth"`
	Gender      string         `gorm:"type:varchar(1)" rql:"filter" json:"gender"`
	Email       string         `gorm:"type:varchar(100)" rql:"filter" json:"email"`
	AddressID   int            `rql:"filter" json:"address_id"`
	Address     Address        `gorm:"foreignKey:AddressID"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

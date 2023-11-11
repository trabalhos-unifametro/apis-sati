package models

import (
	"gorm.io/gorm"
	"time"
)

type Address struct {
	gorm.Model
	ID           uint           `json:"id" gorm:"primaryKey"`
	Street       string         `gorm:"type:varchar" rql:"filter" json:"street"`
	Neighborhood string         `gorm:"type:varchar" rql:"filter" json:"neighborhood"`
	ZipCode      string         `gorm:"type:varchar(9)" rql:"filter" json:"zip_code"`
	Number       string         `gorm:"type:varchar(50)" rql:"filter" json:"number"`
	Complement   string         `gorm:"type:varchar" rql:"filter" json:"complement"`
	CityID       int            `rql:"filter" json:"city_id"`
	City         City           `gorm:"foreignKey:CityID"`
	StateID      int            `rql:"filter" json:"state_id"`
	State        State          `gorm:"foreignKey:StateID"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Address) TableName() string {
	return "address"
}

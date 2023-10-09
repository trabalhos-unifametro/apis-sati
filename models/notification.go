package models

import (
	"gorm.io/gorm"
	"time"
)

type Notification struct {
	gorm.Model
	ID        uint           `json:"id" gorm:"primaryKey"`
	Title     string         `gorm:"type:varchar(100)" rql:"filter" json:"title"`
	Message   string         `gorm:"type:varchar(500)" rql:"filter" json:"message"`
	Type      string         `gorm:"type:varchar(50)" rql:"filter" json:"type"`
	Read      bool           `rql:"filter" json:"read"`
	PatientID int            `rql:"filter" json:"patient_id"`
	Patient   Patient        `gorm:"foreignKey:PatientID"`
	UnitID    int            `rql:"filter" json:"unit_id"`
	Unit      Unit           `gorm:"foreignKey:UnitID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

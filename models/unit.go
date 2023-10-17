package models

import (
	"gorm.io/gorm"
	"time"
)

type Unit struct {
	gorm.Model
	ID                     uint           `json:"id" gorm:"primaryKey"`
	MaxPatientCapacity     int            `json:"max_patient_capacity"`
	CurrentPatientCapacity int            `json:"current_patient_capacity"`
	Name                   string         `json:"name"`
	Status                 string         `json:"status"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

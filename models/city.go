package models

import (
	"gorm.io/gorm"
	"time"
)

type City struct {
	gorm.Model
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `gorm:"type:varchar" rql:"filter" json:"name"`
	StateID   int            `rql:"filter" json:"state_id"`
	State     State          `gorm:"foreignKey:StateID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

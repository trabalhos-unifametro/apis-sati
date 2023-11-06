package models

import (
	"gorm.io/gorm"
	"time"
)

type Status struct {
	gorm.Model
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `gorm:"type:varchar" rql:"filter" json:"name"`
	Description string         `rql:"filter" json:"description"`
	Table       string         `rql:"filter" json:"table"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (Status) TableName() string {
	return "status"
}

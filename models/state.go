package models

import (
	"gorm.io/gorm"
	"time"
)

type State struct {
	gorm.Model
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `gorm:"type:varchar" rql:"filter" json:"name"`
	Abbreviation string         `gorm:"type:varchar(2)" rql:"filter" json:"abbreviation"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

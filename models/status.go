package models

import (
	"apis-sati/database"
	"apis-sati/utils"
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

type ResponseStatus struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (Status) TableName() string {
	return "status"
}

func (s *Status) GetListStatus(tableName string) (error, []ResponseStatus) {
	db := database.OpenConnection()
	var status []ResponseStatus

	err := db.Table("status").
		Select("id, name").
		Where(`"table" ILIKE ?`, tableName).
		Scan(&status).Error

	if err != nil {
		utils.LogMessage{Title: "[MODELS>STATUS] Error on *Status.GetListStatus()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>STATUS] Error on database.CloseConnection(db) > *Status.GetListStatus()", Body: err.Error()}.Error()
	}
	return err, status
}

package models

import (
	"apis-sati/database"
	"apis-sati/utils"
)

func RunMigrations() {
	db := database.OpenConnection()
	if err := db.AutoMigrate(User{}); err != nil {
		utils.LogMessage{Title: "[MIGRATIONS] Error on db.AutoMigrate(User{})", Body: err.Error()}.Error()
	}
	if err := db.AutoMigrate(State{}); err != nil {
		utils.LogMessage{Title: "[MIGRATIONS] Error on db.AutoMigrate(State{})", Body: err.Error()}.Error()
	}
	if err := db.AutoMigrate(City{}); err != nil {
		utils.LogMessage{Title: "[MIGRATIONS] Error on db.AutoMigrate(City{})", Body: err.Error()}.Error()
	}
	if err := db.AutoMigrate(Address{}); err != nil {
		utils.LogMessage{Title: "[MIGRATIONS] Error on db.AutoMigrate(Address{})", Body: err.Error()}.Error()
	}
	if err := db.AutoMigrate(Patient{}); err != nil {
		utils.LogMessage{Title: "[MIGRATIONS] Error on db.AutoMigrate(Patient{})", Body: err.Error()}.Error()
	}
	if err := db.AutoMigrate(Status{}); err != nil {
		utils.LogMessage{Title: "[MIGRATIONS] Error on db.AutoMigrate(Status{})", Body: err.Error()}.Error()
	}
	if err := db.AutoMigrate(Unit{}); err != nil {
		utils.LogMessage{Title: "[MIGRATIONS] Error on db.AutoMigrate(Unit{})", Body: err.Error()}.Error()
	}
	if err := db.AutoMigrate(Notification{}); err != nil {
		utils.LogMessage{Title: "[MIGRATIONS] Error on db.AutoMigrate(Notification{})", Body: err.Error()}.Error()
	}
	if err := db.AutoMigrate(MedicalRecord{}); err != nil {
		utils.LogMessage{Title: "[MIGRATIONS] Error on db.AutoMigrate(MedicalRecord{})", Body: err.Error()}.Error()
	}
}

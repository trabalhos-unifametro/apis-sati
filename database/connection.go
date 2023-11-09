package database

import (
	"apis-sati/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"time"
)

func OpenConnection() *gorm.DB {
	var db *gorm.DB
	str := utils.GodotEnv("CONNECTION_DB")
	database, err := gorm.Open(postgres.Open(str), &gorm.Config{
		NowFunc: func() time.Time {
			ti, _ := time.LoadLocation("UTC")
			return time.Now().In(ti)
		},
	})
	if err != nil {
		utils.LogMessage{Title: "[DATABASE>CONNECTION] Error on OpenConnection()", Body: err.Error()}.Error()
	}

	db = database
	config, _ := db.DB()
	config.SetMaxIdleConns(100)
	config.SetMaxOpenConns(100)
	config.SetConnMaxLifetime(time.Hour)

	return db
}

func CloseConnection(connection *gorm.DB) error {
	db, err := connection.DB()
	if err != nil {
		return err
	}
	if err = db.Close(); err != nil {
		return err
	}

	return nil
}

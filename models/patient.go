package models

import (
	"apis-sati/database"
	"apis-sati/utils"
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

type PatientTotalizatorsDashboard struct {
	Total         int `json:"total"`
	WithinPeriod  int `json:"within_period"`
	OutsidePeriod int `json:"outside_period"`
}

type PatientGraphicDashboard struct {
	PatientsInUnit      int `json:"patients_in_unit"`
	PatientsWaitingUnit int `json:"patients_waiting_unit"`
}

func (p *Patient) TotalizatorsDashboard() (error, PatientTotalizatorsDashboard) {
	db := database.OpenConnection()
	var totalizators PatientTotalizatorsDashboard
	err := db.Table("medical_records mr").
		Select(`COUNT(mr.id) as total,
			SUM(CASE WHEN mr.current_hospitalization_time > mr.expected_hospitalization_time THEN 1 ELSE 0 END) as outside_period,
			SUM(CASE WHEN mr.current_hospitalization_time <= mr.expected_hospitalization_time THEN 1 ELSE 0 END) as within_period`).
		Joins("LEFT JOIN patients p on mr.patient_id = p.id").
		Where("mr.status_id = 1").
		Find(&totalizators).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>PATIENT] Error on *Patient.TotalizatorsDashboard()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>PATIENT] Error on database.CloseConnection(db) > *Patient.TotalizatorsDashboard()", Body: err.Error()}.Error()
	}
	return err, totalizators
}

func (p *Patient) GraphicDashboard() (error, PatientGraphicDashboard) {
	db := database.OpenConnection()
	var graphic PatientGraphicDashboard
	err := db.Table("medical_records mr").
		Select(`
			SUM(CASE WHEN status_id = 1 THEN 1 ELSE 0 END) as patients_in_unit,
			SUM(CASE WHEN status_id = 2 THEN 1 ELSE 0 END) as patients_waiting_unit`).
		Joins("LEFT JOIN patients p on mr.patient_id = p.id").
		Find(&graphic).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>PATIENT] Error on *Patient.GraphicDashboard()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>PATIENT] Error on database.CloseConnection(db) > *Patient.GraphicDashboard()", Body: err.Error()}.Error()
	}
	return err, graphic
}

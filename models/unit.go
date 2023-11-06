package models

import (
	"apis-sati/database"
	"apis-sati/utils"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Unit struct {
	gorm.Model
	ID                     uint           `json:"id" gorm:"primaryKey"`
	MaxPatientCapacity     int            `json:"max_patient_capacity"`
	CurrentPatientCapacity int            `json:"current_patient_capacity"`
	Name                   string         `json:"name"`
	Status                 Status         `gorm:"foreignKey:StatusID"`
	StatusID               int            `rql:"filter" json:"status_id"`
	CreatedAt              time.Time      `json:"created_at"`
	UpdatedAt              time.Time      `json:"updated_at"`
	DeletedAt              gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type UnitTotalizatorsDashboard struct {
	Total         int `json:"total"`
	WithVacancies int `json:"with_vacancies"`
	NoVacancy     int `json:"no_vacancy"`
}

type UnitTotalizators struct {
	QtdMax      int    `json:"qtd_max"`
	QtdPatients int    `json:"qtd_patients"`
	Vacancies   int    `json:"vacancies"`
	UnitName    string `json:"unit_name"`
}

type UnitGraphicDashboard struct {
	WithVacancies int `json:"with_vacancies"`
	NoVacancy     int `json:"no_vacancy"`
}

type ResponseUnit struct {
	ID                     int    `json:"id"`
	Name                   string `json:"name"`
	MaxPatientCapacity     int    `json:"max_patient_capacity"`
	CurrentPatientCapacity int    `json:"current_patient_capacity"`
	NumberOfVacancies      int    `json:"number_of_vacancies"`
	Status                 string `json:"status"`
}

func (u *Unit) TotalizatorsDashboard(unitName, occupation string) (error, UnitTotalizatorsDashboard) {
	db := database.OpenConnection()
	var totalizators UnitTotalizatorsDashboard
	var where string

	if len(strings.TrimSpace(unitName)) > 0 {
		where += fmt.Sprint(`name ILIKE '%`, unitName, `%'`)
	}

	if strings.ToUpper(occupation) == "COM VAGAS" {
		if len(where) > 0 {
			where += " AND "
		}
		where += "current_patient_capacity < max_patient_capacity"
	} else if strings.ToUpper(occupation) == "OCUPADA" {
		if len(where) > 0 {
			where += " AND "
		}
		where += "current_patient_capacity = max_patient_capacity"
	}

	err := db.Table("units").
		Select(`COUNT(*) as total,
			SUM(CASE WHEN current_patient_capacity >= max_patient_capacity THEN 1 ELSE 0 END) as no_vacancy,
			SUM(CASE WHEN current_patient_capacity < max_patient_capacity THEN 1 ELSE 0 END) as with_vacancies`).
		Where(where).
		Scan(&totalizators).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>UNIT] Error on *Unit.TotalizatorsDashboard()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>UNIT] Error on database.CloseConnection(db) > *Unit.TotalizatorsDashboard()", Body: err.Error()}.Error()
	}
	return err, totalizators
}

func (u *Unit) GraphicDashboard() (error, UnitGraphicDashboard) {
	db := database.OpenConnection()
	var graphic UnitGraphicDashboard
	err := db.Table("units").
		Select(`
			SUM(CASE WHEN current_patient_capacity >= max_patient_capacity THEN 1 ELSE 0 END) as no_vacancy,
			SUM(CASE WHEN current_patient_capacity < max_patient_capacity THEN 1 ELSE 0 END) as with_vacancies`).
		Find(&graphic).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>UNIT] Error on *Unit.GraphicDashboard()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>UNIT] Error on database.CloseConnection(db) > *Unit.GraphicDashboard()", Body: err.Error()}.Error()
	}
	return err, graphic
}

func (u *Unit) GetListUnits(unitName, occupation, sortByUnit, sortByOccupation string) (error, []ResponseUnit) {
	db := database.OpenConnection()
	var units []ResponseUnit
	var where string
	var order string

	if len(strings.TrimSpace(unitName)) > 0 {
		where += fmt.Sprint(`u.name ILIKE '%`, unitName, `%'`)
	}

	if strings.ToUpper(occupation) == "COM VAGAS" {
		if len(where) > 0 {
			where += " AND "
		}
		where += "u.current_patient_capacity < u.max_patient_capacity"
	} else if strings.ToUpper(occupation) == "OCUPADA" {
		if len(where) > 0 {
			where += " AND "
		}
		where += "u.current_patient_capacity = u.max_patient_capacity"
	}

	if strings.ToUpper(sortByUnit) == "CRESCENTE" {
		order += "u.name ASC"
	} else if strings.ToUpper(sortByUnit) == "DECRESCENTE" {
		order += "u.name DESC"
	}

	if strings.ToUpper(sortByOccupation) == "CRESCENTE" {
		if len(order) > 0 {
			order += ", "
		}
		order += "number_of_vacancies ASC"
	} else if strings.ToUpper(sortByOccupation) == "DECRESCENTE" {
		if len(order) > 0 {
			order += ", "
		}
		order += "number_of_vacancies DESC"
	}

	err := db.Table("units u").
		Select(`u.id, u.name, u.current_patient_capacity, u.max_patient_capacity,
			(u.max_patient_capacity - u.current_patient_capacity) as number_of_vacancies,
			s.name as status`).
		Joins("LEFT JOIN status s ON s.id = u.status_id").
		Where(where).
		Order(order).
		Scan(&units).Error

	if err != nil {
		utils.LogMessage{Title: "[MODELS>UNIT] Error on *Unit.GetListUnits()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>UNIT] Error on database.CloseConnection(db) > *Unit.GetListUnits()", Body: err.Error()}.Error()
	}
	return err, units
}

func (u *Unit) Totalizators(patientName, situationPatient string) (error, UnitTotalizators) {
	db := database.OpenConnection()
	var totalizators UnitTotalizators
	where := fmt.Sprint("u.id = ", u.ID)

	if len(strings.TrimSpace(patientName)) > 0 {
		where += fmt.Sprint(` AND p.name ILIKE '%`, patientName, `%'`)
	}

	if strings.ToUpper(situationPatient) == "DENTRO DO PERÍODO" {
		where += " AND m.current_hospitalization_time <= m.expected_hospitalization_time"
	} else if strings.ToUpper(situationPatient) == "FORA DO PERÍODO" {
		where += " AND m.current_hospitalization_time > m.expected_hospitalization_time"
	}

	err := db.Table("units u").
		Select(`DISTINCT u.id, 
			 u.name as unit_name, u.current_patient_capacity as qtd_patients, u.max_patient_capacity as qtd_max,
			(u.max_patient_capacity - u.current_patient_capacity) as vacancies`).
		Joins("LEFT JOIN medical_records m ON m.unit_id = u.id").
		Joins("LEFT JOIN patients p ON p.id = m.patient_id").
		Where(where).
		Scan(&totalizators).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>UNIT] Error on *Unit.Totalizators()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>UNIT] Error on database.CloseConnection(db) > *Unit.Totalizators()", Body: err.Error()}.Error()
	}
	return err, totalizators
}

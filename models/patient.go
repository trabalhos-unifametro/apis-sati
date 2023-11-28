package models

import (
	"apis-sati/database"
	"apis-sati/utils"
	"fmt"
	"gorm.io/gorm"
	"strings"
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

type PatientTotalizators struct {
	Total                int `json:"total"`
	HospitalizedPatients int `json:"hospitalized_patients"`
	WaitingPatients      int `json:"waiting_patients"`
}

type PatientGraphicDashboard struct {
	PatientsInUnit      int `json:"patients_in_unit"`
	PatientsWaitingUnit int `json:"patients_waiting_unit"`
}

type ResponsePatient struct {
	PatientID                   int       `json:"patient_id"`
	UnitID                      int       `json:"unit_id"`
	HospitalizationCode         int       `json:"hospitalization_code"`
	PatientName                 string    `json:"patient_name"`
	ExpectedHospitalizationTime time.Time `json:"expected_hospitalization_time"`
	CurrentHospitalizationTime  time.Time `json:"current_hospitalization_time"`
	CreatedAt                   time.Time `json:"created_at"`
	Situation                   string    `json:"situation"`
	SituationID                 int       `json:"situation_id"`
	Status                      string    `json:"status"`
	StatusID                    int       `json:"status_id"`
	MedicalRecordID             int       `json:"medical_record_id"`
	MotherName                  string    `json:"mother_name"`
	Cpf                         string    `json:"cpf"`
	Gender                      string    `json:"gender"`
	UnitName                    string    `json:"unit_name"`
	Street                      string    `json:"street"`
	Neighborhood                string    `json:"neighborhood"`
	Number                      string    `json:"number"`
	City                        string    `json:"city"`
	State                       string    `json:"state"`
	ZipCode                     string    `json:"zip_code"`
	Complement                  string    `json:"complement"`
}

func (p *Patient) TotalizatorsDashboard() (error, PatientTotalizatorsDashboard) {
	db := database.OpenConnection()
	var totalizators PatientTotalizatorsDashboard
	err := db.Table("sati.medical_records mr").
		Select(`COUNT(mr.id) as total,
			SUM(CASE WHEN current_timestamp > mr.expected_hospitalization_time THEN 1 ELSE 0 END) as outside_period,
			SUM(CASE WHEN current_timestamp <= mr.expected_hospitalization_time THEN 1 ELSE 0 END) as within_period`).
		Joins("LEFT JOIN sati.patients p on mr.patient_id = p.id").
		Where("mr.status_id = 1").
		Scan(&totalizators).Error
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
	err := db.Table("sati.medical_records mr").
		Select(`
			SUM(CASE WHEN status_id = 1 THEN 1 ELSE 0 END) as patients_in_unit,
			SUM(CASE WHEN status_id = 2 THEN 1 ELSE 0 END) as patients_waiting_unit`).
		Joins("LEFT JOIN sati.patients p on mr.patient_id = p.id").
		Scan(&graphic).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>PATIENT] Error on *Patient.GraphicDashboard()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>PATIENT] Error on database.CloseConnection(db) > *Patient.GraphicDashboard()", Body: err.Error()}.Error()
	}
	return err, graphic
}

func (p *Patient) Totalizators(patientName, situation string, statusID int) (error, PatientTotalizators) {
	db := database.OpenConnection()
	var totalizators PatientTotalizators
	where := "m.status_id <> 3"

	if len(strings.TrimSpace(patientName)) > 0 {
		where += fmt.Sprint(` AND p.name ILIKE '%`, patientName, `%'`)
	}

	if statusID > 0 {
		where += fmt.Sprint(` AND m.status_id = `, statusID)
	}

	if strings.ToUpper(situation) == "DENTRO DO PERÍODO" {
		where += " AND current_timestamp <= m.expected_hospitalization_time"
	} else if strings.ToUpper(situation) == "FORA DO PERÍODO" {
		where += " AND current_timestamp > m.expected_hospitalization_time"
	}

	err := db.Table("sati.medical_records m").
		Select(`COUNT(m.id) as total,
			SUM(CASE WHEN m.status_id = 1 THEN 1 ELSE 0 END) as hospitalized_patients,
			SUM(CASE WHEN m.status_id = 2 THEN 1 ELSE 0 END) as waiting_patients`).
		Joins("LEFT JOIN sati.patients p on m.patient_id = p.id").
		Where(where).
		Scan(&totalizators).Error
	if err != nil {
		utils.LogMessage{Title: "[MODELS>PATIENT] Error on *Patient.Totalizators()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>PATIENT] Error on database.CloseConnection(db) > *Patient.Totalizators()", Body: err.Error()}.Error()
	}
	return err, totalizators
}

func (p *Patient) GetListPatients(patientName, situationPatient, sortByPatient string, statusID int) (error, []ResponsePatient) {
	db := database.OpenConnection()
	var list []ResponsePatient
	where := fmt.Sprint("m.status_id <> 3")
	var order string

	if len(strings.TrimSpace(patientName)) > 0 {
		where += fmt.Sprint(` AND p.name ILIKE '%`, patientName, `%'`)
	}

	if statusID > 0 {
		where += fmt.Sprint(` AND m.status_id = `, statusID)
	}

	if strings.ToUpper(situationPatient) == "DENTRO DO PERÍODO" {
		where += " AND current_timestamp <= m.expected_hospitalization_time"
	} else if strings.ToUpper(situationPatient) == "FORA DO PERÍODO" {
		where += " AND current_timestamp > m.expected_hospitalization_time"
	}

	if strings.ToUpper(sortByPatient) == "CRESCENTE" {
		order += "p.name ASC"
	} else if strings.ToUpper(sortByPatient) == "DECRESCENTE" {
		order += "p.name DESC"
	}

	err := db.Table("sati.medical_records m").
		Select(`DISTINCT p.id as patient_id, m.created_at,
			u.id as unit_id, m.hospitalization_code, p.name as patient_name, m.expected_hospitalization_time, m.current_hospitalization_time,
			(CASE WHEN current_timestamp > m.expected_hospitalization_time THEN 'Fora do período' ELSE 'Dentro do período' END) as situation,
			(CASE WHEN current_timestamp > m.expected_hospitalization_time THEN 2 ELSE 1 END) as situation_id,
			m.id as medical_record_id, p.mother_name, p.cpf, p.gender, u.name as unit_name, st.name as status, st.id as status_id,
			a.street, a.neighborhood, a.number, a.complement, a.zip_code, c.name as city, s.abbreviation as state`).
		Joins("LEFT JOIN sati.units u ON u.id = m.unit_id").
		Joins("LEFT JOIN sati.patients p ON p.id = m.patient_id").
		Joins("LEFT JOIN sati.status st ON st.id = m.status_id").
		Joins("LEFT JOIN sati.address a ON a.id = p.address_id").
		Joins("LEFT JOIN sati.cities c ON c.id = a.city_id").
		Joins("LEFT JOIN sati.states s ON s.id = a.state_id").
		Where(where).
		Order(order).
		Scan(&list).Error

	if err != nil {
		utils.LogMessage{Title: "[MODELS>PATIENT] Error on *Patient.GetListPatients()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>PATIENT] Error on database.CloseConnection(db) > *Patient.GetListPatients()", Body: err.Error()}.Error()
	}
	return err, list
}

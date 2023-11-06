package models

import (
	"apis-sati/database"
	"apis-sati/utils"
	"gorm.io/gorm"
	"time"
)

type MedicalRecord struct {
	gorm.Model
	ID                          uint           `json:"id" gorm:"primaryKey"`
	Patient                     Patient        `gorm:"foreignKey:PatientID"`
	PatientID                   int            `rql:"filter" json:"patient_id"`
	Unit                        Unit           `gorm:"foreignKey:UnitID"`
	UnitID                      int            `rql:"filter" json:"unit_id"`
	Status                      Status         `gorm:"foreignKey:StatusID"`
	StatusID                    int            `rql:"filter" json:"status_id"`
	ExpectedHospitalizationTime time.Time      `rql:"filter" json:"expected_hospitalization_time"`
	CurrentHospitalizationTime  time.Time      `rql:"filter" json:"current_hospitalization_time"`
	HospitalizationCode         int            `rql:"filter" json:"hospitalization_code"`
	OpeningDate                 time.Time      `rql:"filter" json:"opening_date"`
	CaregiverContact            string         `gorm:"type:varchar(200)" rql:"filter" json:"caregiver_contact"`
	Doctors                     string         `gorm:"type:varchar" rql:"filter" json:"doctors"`
	Schooling                   string         `gorm:"type:varchar(100)" rql:"filter" json:"schooling"`
	Occupation                  string         `gorm:"type:varchar(200)" rql:"filter" json:"occupation"`
	Limitation                  string         `gorm:"type:varchar" rql:"filter" json:"limitation"`
	Allergy                     string         `gorm:"type:varchar" rql:"filter" json:"allergy"`
	CreatedAt                   time.Time      `json:"created_at"`
	UpdatedAt                   time.Time      `json:"updated_at"`
	DeletedAt                   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type MedicalRecordList struct {
	ID                  int `json:"id"`
	MonthOfEntryToUTI   int `json:"month_of_entry_to_uti"`
	MonthOfExitToUTI    int `json:"month_of_exit_to_uti"`
	YearOfExitToUTI     int `json:"year_of_exit_to_uti"`
	TotalVacanciesUnits int `json:"total_vacancies_units"`
}

type MonthlyChart struct {
	Month               int `json:"month"`
	Year                int `json:"year"`
	TotalVacanciesUnits int `json:"total_vacancies_units"`
	TotalVacancies      int `json:"total_vacancies"`
	TotalOccupation     int `json:"total_occupation"`
}

func (m *MedicalRecord) GraphicDashboard() (error, []MonthlyChart) {
	db := database.OpenConnection()
	year := time.Now().Year()
	graphic := []MonthlyChart{
		{Month: 1, Year: year}, {Month: 2, Year: year}, {Month: 3, Year: year}, {Month: 4, Year: year}, {Month: 5, Year: year}, {Month: 6, Year: year},
		{Month: 7, Year: year}, {Month: 8, Year: year}, {Month: 9, Year: year}, {Month: 10, Year: year}, {Month: 11, Year: year}, {Month: 12, Year: year},
	}
	var list []MedicalRecordList
	err := db.Table("medical_records").
		Select(`id, to_char(created_at, 'MM')::int as month_of_entry_to_uti,
			(CASE
					WHEN expected_hospitalization_time > current_hospitalization_time
					THEN to_char(expected_hospitalization_time, 'MM')::int
					ELSE to_char(current_hospitalization_time, 'MM')::int
			END) as month_of_exit_to_uti,
			(CASE
					WHEN expected_hospitalization_time > current_hospitalization_time
					THEN to_char(expected_hospitalization_time, 'YYYY')::int
					ELSE to_char(current_hospitalization_time, 'YYYY')::int
			END) as year_of_exit_to_uti,
			(SELECT SUM(max_patient_capacity) as max_patient_capacity FROM units WHERE to_char(current_date, 'YYYY') = to_char(created_at, 'YYYY')) as total_vacancies_units`).
		Where(`to_char(current_date, 'YYYY') = to_char(created_at, 'YYYY')`).
		Scan(&list).Error

	if err != nil {
		utils.LogMessage{Title: "[MODELS>MEDICAL_RECORD] Error on *MedicalRecord.GraphicDashboard()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>MEDICAL_RECORD] Error on database.CloseConnection(db) > *MedicalRecord.GraphicDashboard()", Body: err.Error()}.Error()
	}

	for i := 0; i < len(graphic); i++ {
		for _, item := range list {
			if graphic[i].Month >= item.MonthOfEntryToUTI &&
				((graphic[i].Month <= item.MonthOfExitToUTI &&
					graphic[i].Year == item.YearOfExitToUTI) ||
					(graphic[i].Month > item.MonthOfExitToUTI &&
						graphic[i].Year < item.YearOfExitToUTI)) {
				graphic[i].TotalVacanciesUnits = item.TotalVacanciesUnits
				graphic[i].TotalOccupation += 1
				graphic[i].TotalVacancies = graphic[i].TotalVacanciesUnits - graphic[i].TotalOccupation
			}
		}
	}

	return err, graphic
}

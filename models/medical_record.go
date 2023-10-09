package models

import (
	"gorm.io/gorm"
	"time"
)

type MedicalRecord struct {
	gorm.Model
	ID                          uint           `json:"id" gorm:"primaryKey"`
	PatientID                   int            `rql:"filter" json:"patient_id"`
	Patient                     Patient        `gorm:"foreignKey:PatientID"`
	UnitID                      int            `rql:"filter" json:"unit_id"`
	Unit                        Unit           `gorm:"foreignKey:UnitID"`
	Status                      string         `gorm:"type:varchar(50)" rql:"filter" json:"status"`
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

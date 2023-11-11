package models

import (
	"apis-sati/database"
	"apis-sati/utils"
	"fmt"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type Notification struct {
	gorm.Model
	ID          uint           `json:"id" gorm:"primaryKey"`
	Title       string         `gorm:"type:varchar(100)" rql:"filter" json:"title"`
	Message     string         `gorm:"type:varchar(500)" rql:"filter" json:"message"`
	Type        string         `gorm:"type:varchar(50)" rql:"filter" json:"type"`
	PatientID   int            `rql:"filter" json:"patient_id"`
	Patient     Patient        `gorm:"foreignKey:PatientID"`
	UnitID      int            `rql:"filter" json:"unit_id"`
	Unit        Unit           `gorm:"foreignKey:UnitID"`
	UsersIDRead []int          `json:"users_id_read" gorm:"type:int[]"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type ResponseNotification struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	Read      bool      `json:"read"`
	PatientID int       `json:"patient_id"`
	UnitID    int       `json:"unit_id"`
	CreatedAt time.Time `json:"created_at"`
}

func (n Notification) GetListNotificationsByUserID(userID int) (error, []ResponseNotification) {
	db := database.OpenConnection()
	var notifications []ResponseNotification
	columns := fmt.Sprint(`id, title, message, created_at, coalesce(patient_id, 0) as patient_id, 
				coalesce(unit_id, 0) as unit_id, type,
   			(CASE WHEN `, userID, ` = ANY(users_id_read) THEN true ELSE false END) as read`)

	err := db.Table("notifications").
		Select(columns).
		Order("created_at DESC").
		Scan(&notifications).Error

	if err != nil {
		utils.LogMessage{Title: "[MODELS>NOTIFICATION] Error on *Notification.GetListNotificationsByUserID()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>NOTIFICATION] Error on database.CloseConnection(db) > *Notification.GetListNotificationsByUserID()", Body: err.Error()}.Error()
	}
	return err, notifications
}

func (n *Notification) ReturnUsersReadNotificationByID() []int {
	db := database.OpenConnection()
	var usersIDRead []int
	err := db.Table("notifications").
		Select("users_id_read").
		Where("id = ?", n.ID).
		Limit(1).
		Scan(&usersIDRead).Error

	if err != nil {
		utils.LogMessage{Title: "[MODELS>NOTIFICATION] Error on *Notification.ReturnUsersReadNotificationByID()", Body: err.Error()}.Error()
	}

	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>NOTIFICATION] Error on database.CloseConnection(db) > *Notification.ReturnUsersReadNotificationByID()", Body: err.Error()}.Error()
	}

	return usersIDRead
}

func (n *Notification) ReadNotificationByUserID(userID int) bool {
	db := database.OpenConnection()
	var success bool
	var strList []string
	list := n.ReturnUsersReadNotificationByID()
	list = append(list, userID)

	for _, num := range list {
		strList = append(strList, strconv.Itoa(num))
	}
	result := strings.Join(strList, ",")

	err := db.Table("notifications").
		Where("id = ?", n.ID).
		Update("users_id_read", fmt.Sprint(`{`, result, `}`)).Error
	if err != nil {
		success = false
	} else {
		success = true
	}
	if err = database.CloseConnection(db); err != nil {
		utils.LogMessage{Title: "[MODELS>NOTIFICATION] Error on database.CloseConnection(db) > *Notification.ReadNotificationByUserID()", Body: err.Error()}.Error()
	}

	return success
}

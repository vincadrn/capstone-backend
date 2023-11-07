package handlers

import (
	// "net/http"
	// "fmt"
	"time"
	"strconv"

	"gorm.io/gorm"

	"github.com/vincadrn/capstone/models"
)

// Send the number of people to db.
func SendNumberofPeople(db *gorm.DB, item []byte) {
	numofPeople, err := strconv.Atoi(string(item))
	if err != nil {
		return
	}
	createdAt := time.Now()
	numofPeopleData := &models.NumberofPeople{Number: numofPeople, CreatedAt: &createdAt}

	db.Create(&numofPeopleData)
}

// Send single bus arrival to db.
func SendBusArrival(db *gorm.DB, item []byte) {
	createdAt := time.Now()
	busArrivalData := &models.BusArrival{IsArrived: true, CreatedAt: &createdAt}

	db.Create(&busArrivalData)
}

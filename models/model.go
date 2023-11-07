package models

import "time"

type Tabler interface {
	TableName() string
}

type NumberofPeople struct {
	ID       			uint
	Number    			int
	CreatedAt 			*time.Time
}

func (NumberofPeople) TableName() string {
	return "number_of_people"
}

type BusArrival struct {
	ID					uint
	IsArrived			bool	
	CreatedAt			*time.Time
}

type FCMClientToken struct {
	ID					uint
	DeviceName			string
	WillBeNotified		bool
	Token				string
}

func (FCMClientToken) TableName() string {
	return "fcm_client_token"
}

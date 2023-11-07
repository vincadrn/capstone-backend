package models

import "time"

type ShowNumberofPeoplePayload struct {
	Number    			int 		`json:"number"`
	CreatedAt 			*time.Time	`json:"created_at"`
}

type ShowBusArrivalPayload struct {
	IsArrived			bool		`json:"is_arrived"`
	CreatedAt			*time.Time	`json:"created_at"`
}

type PostTokenFromClientPayload struct {
	Token				string		`json:"token"`
	DeviceName			string		`json:"device_name"`
	WillBeNotified		string		`json:"will_be_notified"`
}

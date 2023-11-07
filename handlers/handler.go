package handlers

import (
	"context"
	"encoding/json"
	// "fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"firebase.google.com/go/v4/messaging"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"github.com/vincadrn/capstone/models"
	"github.com/vincadrn/capstone/msg"
)

func PostTokenFromClient(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload models.PostTokenFromClientPayload

		body := r.Body
		decoder := json.NewDecoder(body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&payload); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest) + ": invalid payload", http.StatusBadRequest)
			return
		}

		willBeNotified, err := strconv.ParseBool(payload.WillBeNotified)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest) + ": invalid payload", http.StatusBadRequest)
			return
		}

		if err := msg.SetFCMClientToken(db, payload.DeviceName, payload.Token, willBeNotified); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError) + ": cannot set token", http.StatusInternalServerError)
			return
		}
	}
}

func ShowNumberofPeople(item *[]byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		numofPeople, err := strconv.Atoi(string(*item))
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		createdTime := time.Now()
		
		data := models.ShowNumberofPeoplePayload{Number: numofPeople, CreatedAt: &createdTime}
		res, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(res)
	}
}

func ShowBusArrival(item *[]byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		createdTime := time.Now()
		
		data := models.ShowBusArrivalPayload{IsArrived: true, CreatedAt: &createdTime}
		res, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(res)
	}
}

func ListPeople(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		total, err := strconv.Atoi(chi.URLParam(r, "total"))
		if err != nil || total <= 0 {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var people []models.NumberofPeople
		db.Order("created_at desc").Limit(total).Find(&people)

		var data []models.ShowNumberofPeoplePayload
		for _, p := range people {
			data = append(data, models.ShowNumberofPeoplePayload{Number: p.Number, CreatedAt: p.CreatedAt})
		}
		res, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(res)
	}
}

func ListBusArrival(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		total, err := strconv.Atoi(chi.URLParam(r, "total"))
		if err != nil || total <= 0 {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		var busArrival []models.BusArrival
		db.Order("created_at desc").Limit(total).Find(&busArrival)

		var data []models.ShowBusArrivalPayload
		for _, b := range busArrival {
			data = append(data, models.ShowBusArrivalPayload{IsArrived: b.IsArrived, CreatedAt: b.CreatedAt})
		}
		res, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		w.Write(res)
	}
}

/* Test */

func TestSendMessage(db *gorm.DB, client *messaging.Client, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go func() {
			msg.NotifyArrivingBus(db, client, ctx)
		}()

		w.Write([]byte("OK"))
	}
}

func TestSetToken(db *gorm.DB, deviceName string, token string, willBeNotified bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go func() {
			msg.SetFCMClientToken(db, deviceName, token, willBeNotified)
		}()

		w.Write([]byte("OK"))
	}
}

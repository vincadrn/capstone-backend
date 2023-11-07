package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"firebase.google.com/go/v4"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/vincadrn/capstone/handlers"
	"github.com/vincadrn/capstone/models"
	"github.com/vincadrn/capstone/msg"
)

func init() {
	fmt.Println("Initializing ...")
	godotenv.Load(".env")
}

func main() {
	var (
		DB_USER string = os.Getenv("DB_USER")
		DB_PASSWORD string = os.Getenv("DB_PASSWORD")
		DB_HOST string = os.Getenv("DB_HOST")
		DB_PORT string = os.Getenv("DB_PORT")
		DB_DB string = os.Getenv("DB_DB")
		AMQP_USER string = os.Getenv("AMQP_USER")
		AMQP_PASSWORD string = os.Getenv("AMQP_PASSWORD")
		AMQP_HOST string = os.Getenv("AMQP_HOST")
		AMQP_DB string = os.Getenv("AMQP_DB")
	)

	ctx := context.Background()
	firebaseConfig := &firebase.Config{ProjectID: "eng-capstone"}
	app, err := firebase.NewApp(ctx, firebaseConfig)
	if err != nil {
		panic(err)
	}
	clientFCM, err := app.Messaging(ctx)
	if err != nil {
		panic(err)
	}
	
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_DB)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&models.BusArrival{})
	db.AutoMigrate(&models.NumberofPeople{})
	db.AutoMigrate(&models.FCMClientToken{})

	amqpURL := fmt.Sprintf("amqps://%s:%s@%s/%s", AMQP_USER, AMQP_PASSWORD, AMQP_HOST, AMQP_DB)
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Print(err)
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		log.Print(err)
	}
	defer ch.Close()

	args := amqp.Table{"x-message-ttl": 360000, "x-max-length": 255}
	q, err := ch.QueueDeclare("test-queue", true, false, false, false, args)
	if err != nil {
		log.Print(err)
	}
	msgs, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Print(err)
	}

	var packet []byte
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			packet = d.Body
			handlers.SendNumberofPeople(db, d.Body)
		}
	}()

	busQueue, err := ch.QueueDeclare("bus-arrival", true, false, false, false, args)
	if err != nil {
		log.Print(err)
	}
	busMsgs, err := ch.Consume(busQueue.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Print(err)
	}

	var busPacket []byte
	go func() {
		for d := range busMsgs {
			log.Printf("Received an arrival message: %s", d.Body)
			busPacket = d.Body
			handlers.SendNumberofPeople(db, d.Body)
			msg.NotifyArrivingBus(db, clientFCM, ctx)
		}
	}()

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("{\"message\":\"Welcome\"}"))
	})
	r.Post("/api/token", handlers.PostTokenFromClient(db))
	r.Get("/api/people", handlers.ShowNumberofPeople(&packet))
	r.Get("/api/bus", handlers.ShowBusArrival(&busPacket))
	r.Get("/api/data/people/{total}", handlers.ListPeople(db))
	r.Get("/api/data/bus/{total}", handlers.ListBusArrival(db))

	// Test
	// r.Get("/api/test-message", handlers.TestSendMessage(db, clientFCM, ctx))
	// r.Put("/api/test-message", handlers.TestSetToken(db, "test-device", "abc1234def", true))

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("Listening on port %s ...\n", port)
	http.ListenAndServe(":" + port, r)
}

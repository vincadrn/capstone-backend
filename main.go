package main

import (
	"os"
	"fmt"

	"net/http"
	"github.com/go-chi/chi/v5"

	"gorm.io/gorm"
	"gorm.io/driver/postgres"
)

var (
	DB_SERVER string = os.Getenv("DB_SERVER")
	DB_USER string = os.Getenv("DB_USER")
	DB_PASSWORD string = os.Getenv("DB_PASSWORD")
	DB_HOST string = os.Getenv("DB_HOST")
	DB_PORT string = os.Getenv("DB_PORT")
	DB_DB string = os.Getenv("DB_DB")
)

func main() {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_DB)
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("msg: Welcome!"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Listening on port %s ...\n", port)
	http.ListenAndServe(":" + port, r)
}

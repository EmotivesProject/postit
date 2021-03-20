package main

import (
	"log"
	"net/http"
	"os"
	"postit/internal/api"

	"github.com/joho/godotenv"
)

func main() {
	router := api.CreateRouter()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	log.Fatal(http.ListenAndServe(host+":"+port, router))
}

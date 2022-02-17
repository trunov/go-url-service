package main

import (
	"log"
	"net/http"

	"github.com/trunov/go-url-service/internal/app/handlers"
)

func main() {
	http.HandleFunc("/", handlers.ShortenHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
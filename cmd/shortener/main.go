package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/trunov/go-url-service/internal/app/handlers"
)

func main() {
	r := chi.NewRouter()

	r.Route("/api/v1/shortener", func(r chi.Router) {
		r.Post("/", handlers.ShortenHandler)
		r.Get("/{id}", handlers.RedirectHandler)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}

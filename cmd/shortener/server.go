package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/trunov/go-url-service/internal/app/handlers"
	"github.com/trunov/go-url-service/internal/app/storage"
)

func StartServer() {
	urls := make(map[string]string, 10)
	s := storage.NewStorage(urls)

	h := handlers.NewHandlers(s)

	r := chi.NewRouter()
	r.Post("/", h.ShortenHandler)
	r.Get("/{id}", h.RedirectHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}

package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/trunov/go-url-service/internal/app/handlers"
	"github.com/trunov/go-url-service/internal/app/storage"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

// сделать конфиг фолдер и там инициализацию проводить ?!

func StartServer() {
	urls := make(map[string]string, 10)

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.BaseURL, cfg.ServerAddress)

	s := storage.NewStorage(urls)

	h := handlers.NewHandlers(s, cfg.BaseURL)

	r := chi.NewRouter()
	r.Post("/", h.ShortenHandler)
	r.Post("/api/shorten", h.NewShortenHandler)
	r.Get("/{id}", h.RedirectHandler)

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}

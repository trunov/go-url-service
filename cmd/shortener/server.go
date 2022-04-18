package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/trunov/go-url-service/internal/app/file"
	"github.com/trunov/go-url-service/internal/app/handlers"
	"github.com/trunov/go-url-service/internal/app/storage"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStorage   string `env:"FILE_STORAGE"`
}

func StartServer() {
	urls := make(map[string]string, 10)

	var cfg Config
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := file.NewConsumer(cfg.FileStorage)
	if err == nil {
		links, err := consumer.ReadLink()
		if err != nil {
			log.Fatal(err)
		}

		for _, link := range links {
			urls[link.Id] = link.URL
		}

		defer consumer.Close()
	}

	s := storage.NewStorage(urls, cfg.FileStorage)

	h := handlers.NewHandlers(s, cfg.BaseURL)

	r := chi.NewRouter()
	r.Post("/", h.ShortenHandler)
	r.Post("/api/shorten", h.NewShortenHandler)
	r.Get("/{id}", h.RedirectHandler)

	log.Fatal(http.ListenAndServe(cfg.ServerAddress, r))
}

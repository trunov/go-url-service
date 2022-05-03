package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/trunov/go-url-service/internal/app/file"
	"github.com/trunov/go-url-service/internal/app/handlers"
	"github.com/trunov/go-url-service/internal/app/middlewares"
	"github.com/trunov/go-url-service/internal/app/storage"
)

func StartServer(cfg Config) {
	urls := make(map[string]string, 10)

	if cfg.fileStorage != "" {
		consumer, err := file.NewConsumer(cfg.fileStorage)
		if err == nil {
			links, err := consumer.ReadLink()
			if err != nil {
				log.Fatal(err)
			}

			for _, link := range links {
				urls[link.ID] = link.URL
			}

			defer consumer.Close()
		}
	}

	s := storage.NewStorage(urls, cfg.fileStorage)

	h := handlers.NewHandlers(s, cfg.baseURL)

	r := chi.NewRouter()

	r.Use(middlewares.GzipHandle)
	r.Use(middlewares.DecompressHandle)

	r.Post("/", h.ShortenHandler)
	r.Post("/api/shorten", h.NewShortenHandler)
	r.Get("/{id}", h.RedirectHandler)

	fmt.Println("server address " + cfg.serverAddress)

	errServer := http.ListenAndServe(cfg.serverAddress, r)

	if errServer != nil {
		log.Println(errServer)
	}
}

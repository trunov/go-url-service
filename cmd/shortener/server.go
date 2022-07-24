package main

import (
	"crypto/aes"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx"

	"github.com/trunov/go-url-service/internal/app/encryption"
	"github.com/trunov/go-url-service/internal/app/file"
	"github.com/trunov/go-url-service/internal/app/handlers"
	"github.com/trunov/go-url-service/internal/app/middlewares"
	"github.com/trunov/go-url-service/internal/app/storage"
)

func StartServer(cfg Config) {
	urls := make(map[string]string, 10)

	// make a sepparate function which will take cfg.fileStorage and urls as arguments if database dsn is not provided
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

	dbConfig, err := pgx.ParseConnectionString(cfg.databaseDSN)
	if err != nil {
		log.Println(err)
	}

	var conn *pgx.Conn

	if cfg.databaseDSN != "" {
		var err error
		conn, err = pgx.Connect(dbConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			os.Exit(1)
		}
		defer conn.Close()

		if _, err := conn.Exec(
			"CREATE TABLE IF NOT EXISTS shortener (id INT PRIMARY KEY, url VARCHAR(256), shortener_key VARCHAR(256), user_id VARCHAR(256), created_at TIMESTAMP default CURRENT_TIMESTAMP)"); err != nil {
			log.Fatal(err)
		}
	}

	key, err := encryption.GenerateRandom(2 * aes.BlockSize)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	s := storage.NewStorage(urls, cfg.fileStorage)

	h := handlers.NewHandlers(s, cfg.baseURL, conn)

	r := chi.NewRouter()

	r.Use(middlewares.GzipHandle)
	r.Use(middlewares.DecompressHandle)
	r.Use(middlewares.CookieMiddleware(key))

	r.Post("/", h.ShortenHandler)
	r.Post("/api/shorten", h.NewShortenHandler)
	r.Get("/{id}", h.RedirectHandler)
	r.Get("/ping", h.PingDbHandler)

	fmt.Println("server address " + cfg.serverAddress)

	errServer := http.ListenAndServe(cfg.serverAddress, r)

	if errServer != nil {
		log.Println(errServer)
	}
}

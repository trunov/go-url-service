package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/trunov/go-url-service/internal/app/file"
	"github.com/trunov/go-url-service/internal/app/handlers"
	"github.com/trunov/go-url-service/internal/app/middlewares"
	"github.com/trunov/go-url-service/internal/app/storage"
)

// make config struct

const serverAddressDefault string = "localhost:8080"
const baseURLDefault string = "http://localhost:8080"

// make config internal
var (
	fileStorage, serverAddress, baseURL string
)

func init() {
	flag.StringVar(&baseURL, "b", baseURLDefault, "BASE_URL")

	if bu, flgBu := os.LookupEnv("BASE_URL"); flgBu {
		baseURL = bu
	}

	flag.StringVar(&serverAddress, "a", serverAddressDefault, "SERVER_ADDRESS")

	if sa, flgSa := os.LookupEnv("SERVER_ADDRESS"); flgSa {
		serverAddress = sa
	}

	flag.StringVar(&fileStorage, "f", "", "FILE_STORAGE_PATH - путь до файла с сокращёнными URL")

	if u, flgFs := os.LookupEnv("FILE_STORAGE_PATH"); flgFs {
		fileStorage = u
	}
}

func StartServer() {
	fmt.Println("start server")
	defer fmt.Println("server stopped")

	urls := make(map[string]string, 10)

	flag.Parse()

	fmt.Printf(`base url = "%s"
	server address = "%s"
	file path = "%s"
	`, baseURL, serverAddress, fileStorage)

	if fileStorage != "" {
		consumer, err := file.NewConsumer(fileStorage)
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

	s := storage.NewStorage(urls, fileStorage)

	h := handlers.NewHandlers(s, baseURL)

	r := chi.NewRouter()

	r.Use(middlewares.GzipHandle)

	r.Post("/", h.ShortenHandler)
	r.Post("/api/shorten", h.NewShortenHandler)
	r.Get("/{id}", h.RedirectHandler)

	fmt.Println("server address " + serverAddress)

	errServer := http.ListenAndServe(serverAddress, r)

	if errServer != nil {
		log.Println(errServer)
	}
}

// создать явный http сервер дефер сервер шат даун

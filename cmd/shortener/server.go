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
	"github.com/trunov/go-url-service/internal/app/storage"
)

const serverAddressDefault string = "localhost:8080"
const baseURLDefault string = "http://localhost:8080" 

var (fileStorage, serverAddress, baseUrl *string)

func init() {
	baseUrl = flag.String("b", "", "BASE_URL")

	if bu, flgBu := os.LookupEnv("BASE_URL"); flgBu {
		*baseUrl = bu
	} else if *baseUrl == "" {
		*baseUrl = baseURLDefault
	}

	serverAddress = flag.String("a", "", "SERVER_ADDRESS")

	if sa, flgSa := os.LookupEnv("SERVER_ADDRESS"); flgSa {
		*serverAddress = sa
	} else if *serverAddress == "" {
		*serverAddress = serverAddressDefault
	}

	fileStorage = flag.String("f", "", "FILE_STORAGE_PATH - путь до файла с сокращёнными URL")

	if u, flgFs := os.LookupEnv("FILE_STORAGE_PATH"); flgFs {
		*fileStorage = u
	}
}

func StartServer() {
	urls := make(map[string]string, 10)

	flag.Parse()

	fmt.Println(*baseUrl)

	consumer, err := file.NewConsumer(*fileStorage)
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

	s := storage.NewStorage(urls, *fileStorage)

	h := handlers.NewHandlers(s, *baseUrl)

	r := chi.NewRouter()
	r.Post("/", h.ShortenHandler)
	r.Post("/api/shorten", h.NewShortenHandler)
	r.Get("/{id}", h.RedirectHandler)

	fmt.Println("server address " + *serverAddress)

	log.Fatal(http.ListenAndServe(*serverAddress, r))
}

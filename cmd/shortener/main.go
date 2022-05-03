package main

import (
	"flag"
	"os"
)

type Config struct {
	fileStorage   string
	serverAddress string
	baseURL       string
}

func flagSetup() Config {
	var cfg Config

	const serverAddressDefault string = "localhost:8080"
	const baseURLDefault string = "http://localhost:8080"

	flag.StringVar(&cfg.baseURL, "b", baseURLDefault, "BASE_URL")

	if bu, flgBu := os.LookupEnv("BASE_URL"); flgBu {
		cfg.baseURL = bu
	}

	flag.StringVar(&cfg.serverAddress, "a", serverAddressDefault, "SERVER_ADDRESS")

	if sa, flgSa := os.LookupEnv("SERVER_ADDRESS"); flgSa {
		cfg.serverAddress = sa
	}

	flag.StringVar(&cfg.fileStorage, "f", "", "FILE_STORAGE_PATH - путь до файла с сокращёнными URL")

	if u, flgFs := os.LookupEnv("FILE_STORAGE_PATH"); flgFs {
		cfg.fileStorage = u
	}

	return cfg
}

func main() {
	cfg := flagSetup()
	StartServer(cfg)
}

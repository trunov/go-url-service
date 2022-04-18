package storage

import (
	"fmt"
	"log"

	"github.com/trunov/go-url-service/internal/app/file"
)

type urls map[string]string

type Storage struct {
	urls     urls
	fileName string
}

func NewStorage(urls urls, fileName string) *Storage {
	return &Storage{
		urls:     urls,
		fileName: fileName,
	}
}

type Storager interface {
	Get(id string) (string, error)
	Add(id, url string)
	GetAll() urls
}

func (s *Storage) Get(id string) (string, error) {
	value, ok := s.urls[id]
	if !ok {
		return "", fmt.Errorf("value %s not found", id)
	}

	return value, nil
}

func (s *Storage) Add(id, url string) {
	s.urls[id] = url

	producer, err := file.NewProducer(s.fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer producer.Close()

	link := &file.Link{
		Id:  id,
		URL: url,
	}

	writeErr := producer.WriteLink(link)
	if writeErr != nil {
		log.Fatal(err)
	}
}

func (s *Storage) GetAll() urls {
	return s.urls
}

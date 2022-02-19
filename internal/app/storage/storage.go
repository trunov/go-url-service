package storage

import (
	"fmt"
)
type urls map[string]string

type Storage struct {
	urls urls
}

func NewStorage(urls urls) *Storage {
	return &Storage{
		urls: urls,
	}
}

// как использовать этот интерфейс в дальнейшем ?
// на данном этапе storage используется в хендлерах с иницилизацией контруктора NewStorage
// 
type Storager interface {
	Get(id string) (string, error)
	Add(id, url string)
	GetAll() (urls)
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
}

func (s *Storage) GetAll() (urls) {
	return s.urls
}
package file

import (
	"bufio"
	"encoding/json"
	"os"
)

type Link struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

type producer struct {
	file    *os.File
	encoder *json.Encoder
}

func NewProducer(filename string) (*producer, error) {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return nil, err
	}

	return &producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *producer) WriteLink(link *Link) error {
	return p.encoder.Encode(&link)
}

func (p *producer) Close() error {
	return p.file.Close()
}

type consumer struct {
	file    *os.File
	scanner *bufio.Scanner
}

func NewConsumer(filename string) (*consumer, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}

	return &consumer{
		file:    file,
		scanner: bufio.NewScanner(file),
	}, nil
}

func (c *consumer) Close() error {
	return c.file.Close()
}

func (c *consumer) ReadLink() ([]Link, error) {
	linkList := []Link{}

	link := Link{}

	for c.scanner.Scan() {
		err := json.Unmarshal(c.scanner.Bytes(), &link)
		if err != nil {
			return nil, err
		}
		linkList = append(linkList, link)
	}

	return linkList, nil
}

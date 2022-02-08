package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

var urls = make(map[string]string, 10)

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func generateString() string {
	bytes := make([]byte, 8)
	for i := 0; i < 8; i++ {
		bytes[i] = byte(randInt(97, 122))
	}

	return string(bytes)
}

func mapkey(m map[string]string, value string) (key string, ok bool) {
	for k, v := range m {
		if v == value {
			key = k
			ok = true
			return
		}
	}
	return
}

func BodyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		b, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		key, ok := mapkey(urls, string(b))

		var newlyGeneratedKey, shortcutUrl string
		if !ok {
			newlyGeneratedKey = generateString()
			urls[newlyGeneratedKey] = string(b)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(201)

		if ok {
			shortcutUrl = "http://localhost:8080/" + key
		} else {
			shortcutUrl = "http://localhost:8080/" + newlyGeneratedKey
		}

		w.Write([]byte(shortcutUrl))
	case "GET":
		q := r.URL.Path
		q = strings.TrimLeft(q, "/")

		fmt.Println(urls[q]);
		w.Header().Set("Location", urls[q])
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write(nil)
	}

}

func main() {
	http.HandleFunc("/", BodyHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

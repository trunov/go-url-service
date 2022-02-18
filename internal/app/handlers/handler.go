package handlers

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/trunov/go-url-service/internal/app"
)

var urls = make(map[string]string, 10)

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

func ShortenHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		b, err := io.ReadAll(r.Body)

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		key, ok := mapkey(urls, string(b))

		var newlyGeneratedShortLink, tinyURL string
		if !ok {
			newlyGeneratedShortLink = app.GenerateShortLink() // импортируемые функции должны быть с заглавной буквы
			urls[newlyGeneratedShortLink] = string(b)
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(201)

		if ok {
			tinyURL = "http://localhost:8080/" + key
		} else {
			tinyURL = "http://localhost:8080/" + newlyGeneratedShortLink
		}

		w.Write([]byte(tinyURL))

	case "GET":
		q := r.URL.Path
		q = strings.TrimLeft(q, "/")

		fmt.Println(urls[q])
		w.Header().Set("Location", urls[q])
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write(nil)
	}
}

package handlers

import (
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/trunov/go-url-service/internal/app"
	"github.com/trunov/go-url-service/internal/app/storage"
)

var urls = make(map[string]string, 10)
var s = storage.NewStorage(urls)


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
	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	key, ok := mapkey(s.GetAll(), string(b))

	var newlyGeneratedShortLink, tinyURL string
	if !ok {
		newlyGeneratedShortLink = app.GenerateShortLink()
		s.Add(newlyGeneratedShortLink, string(b))
		// urls[newlyGeneratedShortLink] = string(b)
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(201)

	if ok {
		tinyURL = "http://localhost:8080/" + key
	} else {
		tinyURL = "http://localhost:8080/" + newlyGeneratedShortLink
	}

	w.Write([]byte(tinyURL))
}

func RedirectHandler(rw http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(rw, "id param is missing", http.StatusBadRequest)
		return
	}

	url, _ := s.Get(id)
	if url == "" {
		http.Error(rw, "provided id was not found", http.StatusNotFound)
		return
	}

	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
	rw.Write(nil)
}

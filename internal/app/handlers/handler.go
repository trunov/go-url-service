package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/trunov/go-url-service/internal/app"
	"github.com/trunov/go-url-service/internal/app/storage"
)

type Response struct {
	Result string `json:"result"`
}

type Body struct {
	URL string `json:"url"`
}

const localhost string = "http://localhost:8080/"

type Handlers struct {
	storage storage.Storager
	baseURL string
}

func NewHandlers(storage storage.Storager, baseURL string) *Handlers {
	return &Handlers{
		storage: storage,
		baseURL: baseURL,
	}
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

func (h *Handlers) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	key, ok := mapkey(h.storage.GetAll(), string(b))

	var newlyGeneratedShortLink, tinyURL string
	if !ok {
		newlyGeneratedShortLink = app.GenerateShortLink()
		h.storage.Add(newlyGeneratedShortLink, string(b))
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(201)

	if ok {
		tinyURL = localhost + key
	} else {
		tinyURL = localhost + newlyGeneratedShortLink
	}

	w.Write([]byte(tinyURL))
}

func (h *Handlers) RedirectHandler(rw http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(rw, "id param is missing", http.StatusBadRequest)
		return
	}

	url, _ := h.storage.Get(id)
	if url == "" {
		http.Error(rw, "provided id was not found", http.StatusNotFound)
		return
	}

	rw.Header().Set("Location", url)
	rw.WriteHeader(http.StatusTemporaryRedirect)
	rw.Write(nil)
}

func (h *Handlers) NewShortenHandler(rw http.ResponseWriter, r *http.Request) {
	var b Body

	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	key, ok := mapkey(h.storage.GetAll(), b.URL)

	var newlyGeneratedShortLink, tinyURL string
	if !ok {
		newlyGeneratedShortLink = app.GenerateShortLink()
		h.storage.Add(newlyGeneratedShortLink, b.URL)
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(201)

	if ok {
		tinyURL = h.baseURL + "/" + key
	} else {
		tinyURL = h.baseURL + "/" + newlyGeneratedShortLink
	}

	data := Response{Result: tinyURL}

	json.NewEncoder(rw).Encode(data)
}

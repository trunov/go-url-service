package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/trunov/go-url-service/internal/app/storage"
)

func TestShortenHandler(t *testing.T) {

	// тест который проверяет пост запрос на генерацию ключа в ответе текстом и кодом 201 в последующем проверка того что ключ есть в массиве
	type want struct {
		code        int
		response    string
		contentType string
		url         string
		method      string
	}
	tests := []struct {
		name string
		want want
	}{
		{
			name: "Should return 201 in status code and url key as a 5 characters string in json format",
			want: want{
				code:        201,
				response:    "http://localhost:8080",
				contentType: "application/json; charset=utf-8",
				url:         "https://yourbasic.org/golang/io-reader-interface-explained/",
				method:      http.MethodPost,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := strings.NewReader(`{"url": "` + tt.want.url + `"}`)
			request := httptest.NewRequest(tt.want.method, "/api/shorten", link)

			urls := make(map[string]string, 10)
			s := storage.NewStorage(urls)
			handlers := NewHandlers(s)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(handlers.ShortenHandler)
			h.ServeHTTP(w, request)
			res := w.Result()

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			var result Response

			err := json.NewDecoder(res.Body).Decode(&result)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			defer res.Body.Close()

			require.NoError(t, err)
			assert.Contains(t, result.Result, tt.want.response)
		})
	}
}

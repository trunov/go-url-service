package handlers

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			name: "Should return 201 in status code and generatedShortenKey as a 5 characters string in plain text",
			want: want{
				code:        201,
				response:    "http://localhost:8080",
				contentType: "text/plain; charset=utf-8",
				url:         "https://yourbasic.org/golang/io-reader-interface-explained/",
				method:      http.MethodPost,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			link := strings.NewReader(tt.want.url)
			request := httptest.NewRequest(tt.want.method, "/", link)

			w := httptest.NewRecorder()
			h := http.HandlerFunc(ShortenHandler)
			h.ServeHTTP(w, request)
			res := w.Result()

			log.Println("statusCode", res.StatusCode)

			assert.Equal(t, tt.want.code, res.StatusCode)
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)

			assert.Contains(t, string(resBody), tt.want.response)
		})
	}
}
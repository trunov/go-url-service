package middlewares

import (
	"context"
	"crypto/aes"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/trunov/go-url-service/internal/app/encryption"
)

const cookieName = "user_id"
var ctxName interface{} = "user_id"

// make const for cookie key add into env


func CookieMiddleware(key []byte) func(next http.Handler) http.Handler  {
	return func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key, err := encryption.GenerateRandom(2 * aes.BlockSize) 
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return	
		}

		cookieUserID, _ := r.Cookie(cookieName)
		encryptor := encryption.NewEncryptor(key) 

		if cookieUserID != nil {
			userID, err := encryptor.Decode(cookieUserID.Value)

			fmt.Println(userID)

			if err == nil {
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxName, userID)))
				return
		}
	}
	
	userID := uuid.New().String()
	encoded, err := encryptor.Encode([]byte(userID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(encoded)

	cookie := &http.Cookie{Name: cookieName, Value: encoded, HttpOnly: false}
	http.SetCookie(w, cookie)
	next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxName, userID)))
	})
}
}

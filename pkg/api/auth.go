package api

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
)

func makeToken(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func checkToken(token string, password string) bool {
	expected := makeToken(password)
	return token == expected
}

func auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		pass := os.Getenv("TODO_PASSWORD")

		// если пароль не задан - пропуск
		if pass == "" {
			next(w, r)
			return
		}

		cookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		if !checkToken(cookie.Value, pass) {
			http.Error(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

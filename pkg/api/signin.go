package api

import (
	"encoding/json"
	"net/http"
	"os"
)

type SigninRequest struct {
	Password string `json:"password"`
}

func signinHandler(w http.ResponseWriter, r *http.Request) {
	var req SigninRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")

	// если пароль не задан
	if pass == "" {
		writeJSON(w, map[string]string{"error": "Пароль не задан"})
		return
	}

	if req.Password != pass {
		writeJSON(w, map[string]string{"error": "Неверный пароль"})
		return
	}

	token := makeToken(pass)

	writeJSON(w, map[string]string{
		"token": token,
	})
}

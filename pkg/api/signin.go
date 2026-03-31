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
		writeJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Некорректный JSON",
		})
		return
	}

	pass := os.Getenv("TODO_PASSWORD")

	// пароль не задан на сервере
	if pass == "" {
		writeJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Пароль не задан",
		})
		return
	}

	// неверный пароль
	if req.Password != pass {
		writeJSON(w, http.StatusUnauthorized, map[string]string{
			"error": "Неверный пароль",
		})
		return
	}

	token := makeToken(pass)

	// успех
	writeJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	token, _, err := jwtauth.FromContext(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to read token", err), http.StatusInternalServerError)
		return
	}

	u, ok := token.Get("userId")
	if !ok {
		http.Error(w, fmt.Sprintf("%v\n\nuser info not found", err), http.StatusUnauthorized)
		return
	}

	userID, err := uuid.Parse(u.(string))
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nuser info not found", err), http.StatusUnauthorized)
		return
	}
	balance, err := storage.Stor.GetBalance(context.Background(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to get balance", err), http.StatusUnauthorized)
		return
	}
	body, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to serialize balance", err), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(body))
}

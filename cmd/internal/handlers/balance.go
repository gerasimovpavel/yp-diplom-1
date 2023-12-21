package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"io"
	"net/http"
)

func GetBalance(w http.ResponseWriter, r *http.Request) {
	userID, errorString, status := UserIdFromToken(r)
	if errorString != "" {
		http.Error(w, errorString, status)
		return
	}

	balance, err := storage.Stor.GetBalance(context.Background(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to get balance", err), http.StatusInternalServerError)
		return
	}
	body, err := json.Marshal(balance)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to serialize balance", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(body))
}

package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"io"
	"net/http"
)

func LoadWithdrawals(w http.ResponseWriter, r *http.Request) {
	userID, errorString, status := UserIDFromToken(r)
	if errorString != "" {
		http.Error(w, errorString, status)
		return
	}

	wd, err := storage.Stor.GetWithdrawals(context.Background(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to get withdrawals", err), http.StatusInternalServerError)
		return
	}
	if len(wd) == 0 {
		http.Error(w, "no records", http.StatusNoContent)
		return
	}

	body, err := json.Marshal(wd)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to serialize withdraw", err), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = io.WriteString(w, string(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to write to response writer", err), http.StatusInternalServerError)
		return
	}
}

func SaveWithdraw(w http.ResponseWriter, r *http.Request) {
	userID, errorString, status := UserIDFromToken(r)
	if errorString != "" {
		http.Error(w, errorString, status)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to read body", err), http.StatusInternalServerError)
		return
	}
	r.Body.Close()

	wd := &model.Withdraw{}
	err = json.Unmarshal(body, wd)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to deserialize body", err), http.StatusInternalServerError)
		return
	}
	wd.UserID = userID
	_, err = storage.Stor.SetWithdraw(context.Background(), wd)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to write withdraw", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = io.WriteString(w, string(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to write to response writer", err), http.StatusInternalServerError)
		return
	}
}

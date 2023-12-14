package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
)

func GetWithdrawals(w http.ResponseWriter, r *http.Request) {
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
	io.WriteString(w, string(body))
}

func PostWithdraw(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to read body", err), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	wd := &model.Withdraw{}
	err = json.Unmarshal(body, wd)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to deserialize body", err), http.StatusInternalServerError)
		return
	}
	_, err = storage.Stor.SetWithdraw(context.Background(), wd)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to write withdraw", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(body))
}

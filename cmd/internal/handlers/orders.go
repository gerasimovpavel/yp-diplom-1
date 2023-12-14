package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/accruals"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

func PostOrders(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to read  body", err), http.StatusInternalServerError)
		return
	}

	err = goluhn.Validate(string(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\norder number is invalid", err), http.StatusUnprocessableEntity)
		return
	}

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
	dt := time.Now()
	o := &model.Order{
		Number:     string(body),
		UserID:     userID,
		Status:     "NEW",
		UploadedAt: dt}

	o, err = storage.Stor.SetOrder(context.Background(), o)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to write order", err), http.StatusInternalServerError)
		return
	}

	if !o.UploadedAt.IsZero() {
		if userID != o.UserID {
			http.Error(w, "order added by another user", http.StatusConflict)
			return
		}

		if !dt.Equal(o.UploadedAt) {
			http.Error(w, "order already exist", http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusAccepted)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {

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

	orders, err := storage.Stor.GetOrderByUser(context.Background(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to get order", err), http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		http.Error(w, "", http.StatusNoContent)
		return
	}

	body, err := json.Marshal(orders)
	if err != nil {
		http.Error(w, fmt.Sprintf("%s\n\nне могу сериализовать в json", err.Error()), http.StatusInternalServerError)
		return
	}

	accruals.CheckAccruals(userID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, string(body))
}

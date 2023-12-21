package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/accruals"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/google/uuid"
	"io"
	"net/http"
	"time"
)

func PostOrders(w http.ResponseWriter, r *http.Request) {
	userID, errorString, status := UserIDFromToken(r)
	if errorString != "" {
		http.Error(w, errorString, status)
		return
	}

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

	orderID := uuid.New()

	o := &model.Order{
		OrderID:    orderID,
		Number:     string(body),
		UserID:     userID,
		Status:     "NEW",
		UploadedAt: time.Now().Round(0).Truncate(time.Second)}

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

		if orderID != o.OrderID {

			http.Error(w, "order already exist", http.StatusOK)
			return
		}
	}

	w.WriteHeader(http.StatusAccepted)
}

func GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, errorString, status := UserIDFromToken(r)
	if errorString != "" {
		http.Error(w, errorString, status)
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

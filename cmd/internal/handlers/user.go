package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/jwt"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func PostUserAuth(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "wrong Content-Type", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to read  body", err), http.StatusInternalServerError)
		return
	}

	a := &model.User{}
	err = json.Unmarshal(body, a)

	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to deserialize body", err), http.StatusInternalServerError)
		return
	}

	if r.URL.Path == "/api/user/register" {
		a, err = storage.Stor.CreateUser(context.Background(), a)
		if err != nil {

			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
				http.Error(w, "login already exist", http.StatusConflict)
				return
			}

			http.Error(w, fmt.Sprintf("%v\n\nfailed to create user", err), http.StatusInternalServerError)
			return
		}
	}
	if r.URL.Path == "/api/user/login" {
		a, err = storage.Stor.GetUser(context.Background(), a)

		if err != nil {
			http.Error(w, fmt.Sprintf("%v\n\nfailed to find user/password", err), http.StatusUnauthorized)
			return
		}
	}
	tokenString, err := jwt.CreateToken(a)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed create auth token", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "")
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	w.WriteHeader(http.StatusOK)
}

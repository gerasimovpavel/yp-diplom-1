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

func CheckUserLoginPassword(r *http.Request) (*model.User, int, error) {
	if r.Header.Get("Content-Type") != "application/json" {
		return nil, http.StatusBadRequest, errors.New("wrong Content-Type")
	}

	body, err := io.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("%v\n\nfailed to read  body: ", err)
	}

	user := &model.User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, http.StatusInternalServerError, fmt.Errorf("%v\n\nfailed to deserialize body", err)

	}
	return user, http.StatusOK, nil
}

func Login(w http.ResponseWriter, r *http.Request) {
	user, code, err := CheckUserLoginPassword(r)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	user, err = storage.Stor.GetUser(context.Background(), user)

	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed to find user/password", err), http.StatusUnauthorized)
		return
	}

	tokenString, err := jwt.CreateToken(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed create auth token", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "")
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	w.WriteHeader(http.StatusOK)
}

func Register(w http.ResponseWriter, r *http.Request) {
	user, code, err := CheckUserLoginPassword(r)
	if err != nil {
		http.Error(w, err.Error(), code)
		return
	}

	user, err = storage.Stor.CreateUser(context.Background(), user)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			http.Error(w, "login already exist", http.StatusConflict)
			return
		}

		http.Error(w, fmt.Sprintf("%v\n\nfailed to create user", err), http.StatusInternalServerError)
		return
	}

	tokenString, err := jwt.CreateToken(user)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v\n\nfailed create auth token", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "")
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", tokenString))
	w.WriteHeader(http.StatusOK)
}

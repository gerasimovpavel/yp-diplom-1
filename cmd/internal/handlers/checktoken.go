package handlers

import (
	"fmt"
	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"net/http"
)

func UserIdFromToken(r *http.Request) (userID uuid.UUID, errorString string, statusHTTP int) {
	var uID uuid.UUID
	token, _, err := jwtauth.FromContext(r.Context())

	if token == nil || err != nil {
		return uID, fmt.Sprintf("%v\n\nfailed to read token", err), http.StatusUnauthorized
	}

	u, ok := token.Get("userId")
	if !ok {
		return uID, fmt.Sprintf("%v\n\nuser info not found", err), http.StatusUnauthorized
	}

	uID, err = uuid.Parse(u.(string))
	if err != nil {
		return uID, fmt.Sprintf("%v\n\nfailed to parse userId", err), http.StatusUnauthorized

	}
	return uID, "", http.StatusOK
}

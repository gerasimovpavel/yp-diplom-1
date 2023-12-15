package jwt

import (
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/config"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"net/http"
	"time"
)

func CreateToken(user *model.User) (string, error) {
	tok, err := jwt.NewBuilder().
		Issuer("yp.diplom-1").
		Claim("userId", user.UserID).
		Expiration(time.Now().Add(24 * time.Hour)).
		Build()
	if err != nil {
		return "", err
	}
	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS512, []byte(config.HMACSecret)))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

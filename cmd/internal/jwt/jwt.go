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
		Claim("login", user.Login).
		Expiration(time.Now().Add(24 * time.Hour)).
		Build()

	signed, err := jwt.Sign(tok, jwt.WithKey(jwa.HS512, []byte(config.Options.HMACSecret)))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		if token == nil || jwt.Validate(token) != nil {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		next.ServeHTTP(w, r)
	})
}

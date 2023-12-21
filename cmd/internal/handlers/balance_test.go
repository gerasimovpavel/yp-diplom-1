package handlers

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/config"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/jwt"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/go-chi/jwtauth/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_GetBalanceHandlers(t *testing.T) {
	var err error
	config.Options.DatabaseURI = "host=localhost user=shortener password=shortener dbname=gofermart sslmode=disable"
	gofakeit.Seed(0)
	storage.Stor, err = storage.NewPgStorage()
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name         string
		auth         bool
		wantStatuses []int
		hfunc        http.HandlerFunc
	}{
		{
			"Check Get Balance w/o auth",
			false,
			[]int{http.StatusUnauthorized},
			GetBalance,
		},
		{
			"Check Get Balance with auth",
			true,
			[]int{http.StatusOK, http.StatusNoContent},
			GetBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
			ctx := req.Context()
			tokenString := ""

			if tt.auth {
				user := &model.User{}
				user.UserID = gofakeit.UUID()
				user.Login = gofakeit.Username()
				user.Password = gofakeit.Password(true, true, true, false, false, 8)
				tokenString, err = jwt.CreateToken(user)
				if err != nil {
					panic(err)
				}
			}

			ja := jwtauth.New("HS512", []byte(config.HMACSecret), nil)
			token, err := jwtauth.VerifyToken(ja, tokenString)

			if token != nil {
				ctx = jwtauth.NewContext(ctx, token, err)
				req = req.WithContext(ctx)
			}
			w := httptest.NewRecorder()
			var res *http.Response
			tt.hfunc(w, req)
			res = w.Result()
			defer res.Body.Close()

			if !assert.Contains(t, tt.wantStatuses, res.StatusCode) {
				panic(fmt.Errorf("status expect %v actual %v", tt.wantStatuses, res.StatusCode))
			}
		})
	}

}

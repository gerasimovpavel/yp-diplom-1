package handlers

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/config"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/jwt"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/model"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/storage"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwa"
	jwt2 "github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func CreateWrongKeyToken(user *model.User) (string, error) {
	tok, err := jwt2.NewBuilder().
		Issuer("yp.diplom-1").
		Claim("login", user.UserID).
		Expiration(time.Now().Round(0).Truncate(time.Second).Add(24 * time.Hour)).
		Build()
	if err != nil {
		return "", err
	}
	signed, err := jwt2.Sign(tok, jwt2.WithKey(jwa.HS512, []byte(config.HMACSecret)))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func CreateWrongValueToken(user *model.User) (string, error) {
	tok, err := jwt2.NewBuilder().
		Issuer("yp.diplom-1").
		Claim("userID", strconv.Itoa(gofakeit.Number(6, 10))).
		Expiration(time.Now().Round(0).Truncate(time.Second).Add(24 * time.Hour)).
		Build()
	if err != nil {
		return "", err
	}
	signed, err := jwt2.Sign(tok, jwt2.WithKey(jwa.HS512, []byte(config.HMACSecret)))
	if err != nil {
		return "", err
	}

	return string(signed), nil
}

func Test_GetOrders(t *testing.T) {
	var err error
	config.Options.DatabaseURI = "host=localhost user=shortener password=shortener dbname=gofermart sslmode=disable"
	gofakeit.Seed(0)

	storage.Stor, err = storage.NewPgStorage()
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name              string
		tokenKeyCorrect   bool
		tokenValueCorrect bool
		auth              bool
		userID            string
		login             string
		password          string
		order             string
		wantStatuses      []int
		hfunc             http.HandlerFunc
	}{
		{
			"token read error",
			false,
			false,
			false,
			gofakeit.UUID(),
			gofakeit.Username(),
			gofakeit.Password(true, true, true, false, false, 9),
			strconv.Itoa(gofakeit.CreditCardNumberLuhn()),
			[]int{http.StatusUnauthorized},
			GetOrders,
		},
		{
			"token user info error",
			false,
			true,
			true,
			gofakeit.UUID(),
			gofakeit.Username(),
			gofakeit.Password(true, true, true, false, false, 9),
			strconv.Itoa(gofakeit.Number(3, 10)),
			[]int{http.StatusUnauthorized},
			GetOrders,
		},
		{
			"token failed parse userID",
			true,
			false,
			true,
			gofakeit.UUID(),
			gofakeit.Username(),
			gofakeit.Password(true, true, true, false, false, 9),
			strconv.Itoa(gofakeit.Number(3, 10)),
			[]int{http.StatusUnauthorized},
			GetOrders,
		},
		{
			"get orders from storage",
			true,
			true,
			true,
			gofakeit.UUID(),
			gofakeit.Username(),
			gofakeit.Password(true, true, true, false, false, 9),
			strconv.Itoa(gofakeit.Number(3, 10)),
			[]int{http.StatusNoContent, http.StatusOK},
			GetOrders,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/user/orders", nil)
			ctx := req.Context()

			tokenString := ""

			switch tt.tokenKeyCorrect && tt.tokenValueCorrect {
			case true:
				{
					user, err := storage.Stor.CreateUser(ctx, &model.User{Login: tt.login, Password: tt.password})
					if err != nil {
						panic(err)
					}
					tokenString, err = jwt.CreateToken(user)
					if err != nil {
						panic(err)
					}
				}
			case false:
				{
					switch tt.tokenKeyCorrect {
					case false:
						{
							tokenString, err = CreateWrongKeyToken(
								&model.User{
									UserID:   tt.userID,
									Login:    tt.login,
									Password: tt.password})
							if err != nil {
								panic(err)
							}
						}
					case true:
						{
							tokenString, err = CreateWrongValueToken(&model.User{
								UserID:   tt.userID,
								Login:    tt.login,
								Password: tt.password})
							if err != nil {
								panic(err)
							}
						}
					}
				}
			}

			if tt.auth {
				ja := jwtauth.New("HS512", []byte(config.HMACSecret), nil)

				token, err := jwtauth.VerifyToken(ja, tokenString)
				if token != nil {
					ctx = jwtauth.NewContext(ctx, token, err)
					req = req.WithContext(ctx)
				}
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

func Test_PostOrders(t *testing.T) {
	var err error
	config.Options.DatabaseURI = "host=localhost user=shortener password=shortener dbname=gofermart sslmode=disable"
	gofakeit.Seed(0)

	storage.Stor, err = storage.NewPgStorage()
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name              string
		tokenKeyCorrect   bool
		tokenValueCorrect bool
		auth              bool
		userID            string
		login             string
		password          string
		order             string
		wantStatuses      []int
		hfunc             http.HandlerFunc
	}{
		{
			"token read error",
			false,
			false,
			false,
			gofakeit.UUID(),
			gofakeit.Username(),
			gofakeit.Password(true, true, true, false, false, 9),
			strconv.Itoa(gofakeit.CreditCardNumberLuhn()),
			[]int{http.StatusUnauthorized},
			PostOrders,
		},
		{
			"token user info error",
			false,
			true,
			true,
			gofakeit.UUID(),
			gofakeit.Username(),
			gofakeit.Password(true, true, true, false, false, 9),
			strconv.Itoa(gofakeit.Number(3, 10)),
			[]int{http.StatusUnauthorized},
			PostOrders,
		},
		{
			"token failed parse userID",
			true,
			false,
			true,
			gofakeit.UUID(),
			gofakeit.Username(),
			gofakeit.Password(true, true, true, false, false, 9),
			strconv.Itoa(gofakeit.Number(3, 10)),
			[]int{http.StatusUnauthorized},
			PostOrders,
		},
		{
			"post orders(no luhn) from storage",
			true,
			true,
			true,
			gofakeit.UUID(),
			gofakeit.Username(),
			gofakeit.Password(true, true, true, false, false, 9),
			strconv.Itoa(gofakeit.Number(3, 10)),
			[]int{http.StatusUnprocessableEntity},
			PostOrders,
		},
		{
			"post orders(uhn) from storage",
			true,
			true,
			true,
			gofakeit.UUID(),
			gofakeit.Username(),
			gofakeit.Password(true, true, true, false, false, 9),
			strconv.Itoa(gofakeit.CreditCardNumberLuhn()),
			[]int{http.StatusAccepted},
			PostOrders,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req := httptest.NewRequest(http.MethodGet, "/api/user/orders", strings.NewReader(tt.order))
			ctx := req.Context()

			tokenString := ""

			switch tt.tokenKeyCorrect && tt.tokenValueCorrect {
			case true:
				{
					user, err := storage.Stor.CreateUser(ctx, &model.User{Login: tt.login, Password: tt.password})
					if err != nil {
						panic(err)
					}
					tokenString, err = jwt.CreateToken(user)
					if err != nil {
						panic(err)
					}
				}
			case false:
				{
					switch tt.tokenKeyCorrect {
					case false:
						{
							tokenString, err = CreateWrongKeyToken(&model.User{UserID: tt.userID, Login: tt.login, Password: tt.password})
							if err != nil {
								panic(err)
							}
						}
					case true:
						{
							tokenString, err = CreateWrongValueToken(&model.User{UserID: tt.userID, Login: tt.login, Password: tt.password})
							if err != nil {
								panic(err)
							}
						}
					}
				}
			}

			if tt.auth {
				ja := jwtauth.New("HS512", []byte(config.HMACSecret), nil)

				token, err := jwtauth.VerifyToken(ja, tokenString)
				if token != nil {
					ctx = jwtauth.NewContext(ctx, token, err)
					req = req.WithContext(ctx)
				}
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

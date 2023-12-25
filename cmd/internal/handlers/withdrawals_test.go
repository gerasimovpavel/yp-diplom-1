package handlers

import (
	"context"
	"encoding/json"
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
	"strconv"
	"strings"
	"testing"
	"time"
)

func Test_GetWithdrawals(t *testing.T) {
	var err error
	config.Options.DatabaseURI = "host=localhost user=shortener password=shortener dbname=gofermart sslmode=disable"
	gofakeit.Seed(0)

	storage.Stor, err = storage.NewPgStorage()
	if err != nil {
		panic(err)
	}

	user := &model.User{Login: gofakeit.Username(), Password: gofakeit.Password(
		true,
		true,
		true,
		false,
		false,
		10)}

	user, err = storage.Stor.CreateUser(context.Background(), user)
	if err != nil {
		panic(err)
	}

	tests := []struct {
		name         string
		method       string
		path         string
		wantStatuses []int
		hfunc        http.HandlerFunc
	}{
		{
			"Save Withdrawals",
			http.MethodPost,
			"/api/user/balance/withdraw",
			[]int{http.StatusOK},
			SaveWithdraw,
		},
		{
			"Load Withdrawals",
			http.MethodGet,
			"/api/user/withdrawals",
			[]int{http.StatusOK},
			LoadWithdrawals,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokenString, err := jwt.CreateToken(user)
			ja := jwtauth.New("HS512", []byte(config.HMACSecret), nil)

			token, err := jwtauth.VerifyToken(ja, tokenString)

			if err != nil {
				panic(err)
			}

			now := time.Now()
			wd := &model.Withdraw{
				Order:       strconv.Itoa(gofakeit.CreditCardNumberLuhn()),
				Sum:         float64(gofakeit.Number(100, 900)),
				UserID:      user.UserID,
				ProcessedAt: gofakeit.DateRange(now.Add(time.Minute*60*24*30*-1), now),
			}

			wdJSON, err := json.Marshal(wd)
			if err != nil {
				panic(err)
			}

			body := strings.NewReader(string(wdJSON))

			req := httptest.NewRequest(tt.method, tt.path, body)

			ctx := req.Context()
			ctx = jwtauth.NewContext(ctx, token, err)
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()
			var res *http.Response
			tt.hfunc(w, req)

			res = w.Result()
			res.Body.Close()
			if !assert.Contains(t, tt.wantStatuses, res.StatusCode) {
				panic(fmt.Errorf("status expect %v actual %v", tt.wantStatuses, res.StatusCode))
			}

		})
	}
}

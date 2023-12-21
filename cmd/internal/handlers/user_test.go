package handlers

import (
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_PostUser(t *testing.T) {
	login := gofakeit.Username()
	password := gofakeit.Password(true, true, true, false, false, 12)
	tests := []struct {
		name         string
		path         string
		contentType  string
		body         string
		wantStatuses []int
		hfunc        http.HandlerFunc
	}{
		{
			"wrong content-type",
			"/api/user/register",
			"text/plain",
			`{"login":"","password":""}`,
			[]int{http.StatusBadRequest},
			PostUserAuth,
		},
		{
			"failed deserialize",
			"/api/user/register",
			"application/json",
			`{login":"","password":""}`,
			[]int{http.StatusInternalServerError},
			PostUserAuth,
		},
		{
			"login w/o register",
			"/api/user/login",
			"application/json",
			fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password),
			[]int{http.StatusUnauthorized},
			PostUserAuth,
		},
		{
			"register",
			"/api/user/register",
			"application/json",
			fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password),
			[]int{http.StatusOK},
			PostUserAuth,
		},
		{
			"register after register",
			"/api/user/register",
			"application/json",
			fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password),
			[]int{http.StatusConflict},
			PostUserAuth,
		},
		{
			"login after register",
			"/api/user/login",
			"application/json",
			fmt.Sprintf(`{"login":"%s","password":"%s"}`, login, password),
			[]int{http.StatusOK},
			PostUserAuth,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", tt.contentType)
			w := httptest.NewRecorder()
			var res *http.Response
			tt.hfunc(w, req)
			res = w.Result()

			if !assert.Contains(t, tt.wantStatuses, res.StatusCode) {
				panic(fmt.Errorf("status expect %v actual %v", tt.wantStatuses, res.StatusCode))
			}
		})
	}
}

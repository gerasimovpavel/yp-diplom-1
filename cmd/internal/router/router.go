package router

import (
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/config"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/handlers"
	"github.com/gerasimovpavel/yp-diplom-1/cmd/internal/jwt"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
)

var tokenAuth *jwtauth.JWTAuth

func New() chi.Router {
	tokenAuth = jwtauth.New("HS512", []byte(config.HMACSecret), nil)

	r := chi.NewRouter()
	r.Use(
		chimw.Logger,
		chimw.Compress(5),
	)
	r.Route("/", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Route("/user", func(r chi.Router) {
				r.Post("/register", handlers.Register)
				r.Post("/login", handlers.Login)
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwt.Authenticator)
					r.Route("/orders", func(r chi.Router) {
						r.Post("/", handlers.SaveOrders)
						r.Get("/", handlers.LoadOrders)
					})
					r.Route("/balance", func(r chi.Router) {
						r.Get("/", handlers.LoadBalance)
						r.Post("/withdraw", handlers.SaveWithdraw)
					})
					r.Get("/withdrawals", handlers.LoadWithdrawals)
				})
			})
		})

	})
	return r
}

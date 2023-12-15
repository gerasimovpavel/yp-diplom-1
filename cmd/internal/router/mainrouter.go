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

func MainRouter() chi.Router {
	tokenAuth = jwtauth.New("HS512", []byte(config.HMACSecret), nil)

	r := chi.NewRouter()
	r.Use(
		//chizap.New(logger.Logger, &chizap.Opts{
		//	WithReferer:   true,
		//	WithUserAgent: true,
		//}),
		chimw.Logger,
		chimw.Compress(5),
	)
	r.Route("/", func(r chi.Router) {
		r.Route("/api", func(r chi.Router) {
			r.Route("/user", func(r chi.Router) {
				r.Post("/register", handlers.PostUserAuth)
				r.Post("/login", handlers.PostUserAuth)
				r.Group(func(r chi.Router) {
					r.Use(jwtauth.Verifier(tokenAuth))
					r.Use(jwt.Authenticator)
					r.Route("/orders", func(r chi.Router) {
						r.Post("/", handlers.PostOrders)
						r.Get("/", handlers.GetOrders)
					})
					r.Route("/balance", func(r chi.Router) {
						r.Get("/", handlers.GetBalance)
						r.Put("/", handlers.UpdateBalance)
						r.Post("/withdraw", handlers.PostWithdraw)
					})
					r.Get("/withdrawals", handlers.GetWithdrawals)
				})
			})
		})

	})
	return r
}

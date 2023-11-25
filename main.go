package main

import (
	"net/http"
	"strconv"

	"github.com/Hunnnn77/hello/controller"
	"github.com/Hunnnn77/hello/db"
	md "github.com/Hunnnn77/hello/middleware"
	"github.com/Hunnnn77/hello/model"
	"github.com/Hunnnn77/hello/response"
	"github.com/Hunnnn77/hello/util"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

func init() {
	db.Initialize()
	util.JwtAuth = jwtauth.New("HS256", []byte(util.ByField("TOKEN_SECRET")), nil)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:4321"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-JwtAuth"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           3600,
	}))

	r.Use(md.SetHeader)

	r.Group(func(r chi.Router) {
		r.Route("/", func(r chi.Router) {
			r.Post("/login", controller.HandleLogin)
			r.Post("/signup", controller.HandleSignup)
			r.With(md.Authorization).Post("/logout", controller.HandleLogout)
		})

	})

	r.Group(func(r chi.Router) {
		r.Use(md.Authorization)
		r.Route("/auth", func(r chi.Router) {
			r.Get("/", HandleAuth)
		})
	})

	err := http.ListenAndServe("localhost:"+strconv.Itoa(util.PORT), r)
	if err != nil {
		panic(err)
	}
}

func HandleAuth(w http.ResponseWriter, r *http.Request) {
	_, ctx := util.HandleContext(r, nil)
	if tokenString, ok := ctx.(string); !ok {
		response.ThrowError(w, model.HttpError{
			Code: http.StatusUnauthorized,
		})
	} else {
		response.ThrowOk[bool](w, model.HttpOk[bool, string]{
			Ok: true,
			At: &tokenString,
		})
	}
}

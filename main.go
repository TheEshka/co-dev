package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/misgorod/co-dev/auth"
	"github.com/misgorod/co-dev/db"
	"github.com/misgorod/co-dev/middlewares"
	"github.com/misgorod/co-dev/post"
)

func main() {
	client, err := db.Connect()
	if err != nil {
		panic(err)
	}
	authHandler := auth.AuthHandler{
		Client: client,
	}
	postHandler := post.PostHandler{
		Client: client,
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)
	r.Route("/api", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		r.Route("/posts", func(r chi.Router) {
			r.Get("/", postHandler.GetAll)
			r.Group(func(r chi.Router) {
				r.Use(middlewares.Authenticate)
				r.Post("/", postHandler.Post)
			})
			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", postHandler.Get)
				r.Route("/members", func(r chi.Router) {
					r.Use(middlewares.Authenticate)
					r.Post("/", postHandler.MemberPost)
					r.Delete("/", postHandler.MemberDelete)
				})
			})
		})
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}

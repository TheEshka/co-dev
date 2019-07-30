package main

import (
	"github.com/misgorod/co-dev/common"
	"github.com/misgorod/co-dev/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/misgorod/co-dev/middlewares"
	"gopkg.in/go-playground/validator.v9"
)

func main() {
	client, err := common.Connect()
	if err != nil {
		panic(err)
	}
	authHandler := handlers.AuthHandler{
		Client:   client,
		Validate: validator.New(),
	}
	postHandler := handlers.PostsHandler{
		Client:   client,
		Validate: validator.New(),
	}
	usersHandler := handlers.UsersHandler{
		Client: client,
	}
	imagesHandler := handlers.ImagesHandler{
		Client:   client,
	}
	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.Logger, middleware.Recoverer)
	r.Route("/api", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		r.Route("/users", func(r chi.Router) {
			r.Use(middlewares.Authenticate)
			r.Get("/{id}", usersHandler.Get)
			r.Put("/", usersHandler.Put)
			r.Post("/{id}/image", usersHandler.PostImage)
		})

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
					r.Post("/", postHandler.MembersPost)
					r.Delete("/", postHandler.MembersDelete)
					r.Route("/{memberId}", func(r chi.Router) {
						r.Put("/", postHandler.MemberPut)
						r.Delete("/", postHandler.MemberDelete)
					})
				})
				r.Route("/image", func(r chi.Router) {
					r.Use(middlewares.Authenticate)
					r.Post("/", postHandler.PostImage)
				})
			})
		})

		r.Get("/image/{id}", imagesHandler.GetImage)
	})
	log.Fatal(http.ListenAndServe(":8080", r))
}

package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	mux.Get("/", app.GetInitial)
	mux.Get("/book/{id}", app.GetBookById)
	mux.Get("/book", app.ListBooksHandle)
	mux.Post("/book", app.IsAuthorized(app.SaveBookHandle))
	mux.Post("/opinion", app.IsAuthorized(app.SaveOpinionHandle))
	mux.Post("/user", app.SaveUserHandle)
	mux.Post("/user/login", app.LoginUserHandle)
	return mux
}

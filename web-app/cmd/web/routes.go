package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	// register middleware (global)
	mux.Use(middleware.Recoverer)    // middleware para manejar panics
	mux.Use(app.AddIPToContext)      // middleware para agregar IP al contexto
	mux.Use(app.Session.LoadAndSave) // middleware para manejar sesiones

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(app.auth) // proteger rutas debajo de /user es un ejemplo de cómo aplicar middleware a un grupo específico de rutas
		mux.Get("/profile", app.Profile)
	})

	// register routes
	mux.Get("/", app.Home)
	mux.Post("/login", app.Login)

	// static assets
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer)) //registro de ruta template

	return mux
}

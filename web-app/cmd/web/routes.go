package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// routes le pertenece a application por medio del receiver
func (app *application) routes() http.Handler {
	mux := chi.NewRouter() // mux es un enrutador

	// register middleware (global)
	mux.Use(middleware.Recoverer)    // middleware para manejar panics
	mux.Use(app.AddIPToContext)      // middleware para agregar IP al contexto
	mux.Use(app.Session.LoadAndSave) // middleware para manejar sesiones

	mux.Route("/user", func(mux chi.Router) {
		mux.Use(app.auth) // proteger rutas debajo de /user
		mux.Get("/profile", app.Profile)
	})

	// register routes
	mux.Get("/", app.Home)
	mux.Post("/login", app.Login)

	// static assets
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

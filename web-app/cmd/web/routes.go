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

	// register routes
	mux.Get("/", app.Home)
	mux.Post("/login", app.Login)

	// sstatic assets
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

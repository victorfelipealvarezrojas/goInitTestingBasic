package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_application_routes(t *testing.T) {

	type Route struct {
		route  string
		method string
	}

	var registered = []Route{}

	registered = append(registered, Route{
		route:  "/",
		method: "GET",
	})

	registered = append(registered, Route{
		route:  "/static/*",
		method: "GET",
	})

	mux := app.routes()

	chiRoutes := mux.(chi.Routes) // Convierte a chi.Routes para poder inspeccionar

	for _, route := range registered {
		if !routeExists(route.route, route.method, chiRoutes) {
			t.Errorf("The route %s with method %s is not registered", route.route, route.method)
		}
	}
}

func routeExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	found := false

	//  chi.Walk:: recorre las rutas registradas. Para cada una, ejecuta la función anonima
	_ = chi.Walk(
		chiRoutes,
		func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
			// ...func  anónima recibe los detalles de CADA ruta registrada
			// method = el método de ESA ruta (GET, POST, DELETE, etc)
			// route = el path de ESA ruta ("/", "/users", etc)
			// handler = el manejador de ESA ruta
			// middlewares = middlewares de ESA ruta

			if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
				found = true
			}
			return nil
		})

	return found
}

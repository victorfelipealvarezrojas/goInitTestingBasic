package main

import (
	"net/http"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_application_routes(t *testing.T) {

	registered := []struct {
		route  string
		method string
	}{
		{"/", "GET"},
		{"/static/*", "GET"}, // necesita ser registrada
	}

	mux := app.routes()

	chiRoutes := mux.(chi.Routes) // Convierte a chi.Routes Porque app.routes() retorna http.Handler  la vista reducida que solo expone ServeHTTP.
	//Para inspeccionar las rutas registradas necesitas chi.Routes — que expone métodos como Walk. Pero http.Handler no los tiene.

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

package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
)

var app application

// se ejecuta antes de todos los tests
func TestMain(m *testing.M) {
	// pathToTemplates define la ruta a los templates HTML.
	// Durante los tests el working directory es el directorio del paquete (cmd/web/)
	// por lo que la ruta debe ser relativa a ese directorio, no a la raíz del proyecto.
	// Para tests se sobreescribe con "./testdata/" en Test_application_render_bad_template.
	pathToTemplates = "./../../templates/"

	app.Session = getSession()

	app.DB = &dbrepo.TestDBRepo{}

	os.Exit(m.Run()) // ejecutar los tests
}

package main

import (
	"os"
	"testing"
)

var app application

// se ejecuta antes de todos los tests
func TestMain(m *testing.M) {

	pathToTemplates = "./../../templates/"

	app.Session = getSession()

	os.Exit(m.Run()) // ejecutar los tests
}

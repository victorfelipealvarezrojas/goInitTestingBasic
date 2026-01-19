package main

import (
	"os"
	"testing"
	"webapp/pkg/repository/dbrepo"
)

var app application

// se ejecuta antes de todos los tests
func TestMain(m *testing.M) {

	pathToTemplates = "./../../templates/"

	app.Session = getSession()

	app.DB = &dbrepo.TestDBRepo{}

	os.Exit(m.Run()) // ejecutar los tests
}

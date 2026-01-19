package main

import (
	"encoding/gob"
	"flag"
	"log"
	"net/http"
	"webapp/pkg/data"
	"webapp/pkg/repository"
	"webapp/pkg/repository/dbrepo"

	"github.com/alexedwards/scs/v2"
)

type application struct {
	DSN     string
	DB      repository.DataBaseRepository
	Session *scs.SessionManager
}

func main() {

	gob.Register(data.User{})

	app := &application{}

	flag.StringVar(&app.DSN, "dsn", "postgresql://postgres:postgres@localhost:5432/users?sslmode=disable&timezone=UTC", "CNNECTION STRING")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	app.DB = &dbrepo.PostgresDBRepo{DB: conn}

	app.Session = getSession()

	log.Println("Listener Server on port 8080....")

	err = http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

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

	//NOTE: Para almacenar tipos personalizados en la sesión, como data.User, es necesario registrar el tipo con gob.Register. Esto permite que el sistema de codificación de sesiones pueda serializar y deserializar correctamente los valores de ese tipo.
	gob.Register(data.User{})

	app := &application{} // received

	flag.StringVar(&app.DSN, "dsn", "postgresql://postgres:postgres@localhost:5432/users?sslmode=disable&timezone=UTC", "CNNECTION STRING")
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}

	app.DB = &dbrepo.PostgresDBRepo{DataBaseCnn: conn}
	defer conn.Close()

	app.Session = getSession() // para que estio funciuone en los handlers, se asigna a la aplicación antes de llamar a app.routes() que es donde se usa el middleware de sesiones Mux.Use(app.Session.LoadAndSave)

	log.Println("Listener Server on port 8080....")

	err = http.ListenAndServe(":8080", app.routes())
	if err != nil {
		log.Fatal(err)
	}
}

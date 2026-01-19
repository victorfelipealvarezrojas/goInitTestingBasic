package main

import (
	"database/sql"
	"log"

	// pgconn - Driver de bajo nivel para protocolo PostgreSQL
	// Maneja comunicación directa con el servidor (autenticación, envío de comandos)
	_ "github.com/jackc/pgconn"

	// pgx/v4 - Wrapper sobre pgconn que proporciona interfaz más cómoda
	// Ofrece prepared statements, query builders, scanning automático de resultados
	_ "github.com/jackc/pgx/v4"

	// pgx/v4/stdlib - Adapter que convierte pgx al estándar database/sql de Go
	// Permite usar pgx con ORMs y librerías que esperan la interfaz estándar de Go
	_ "github.com/jackc/pgx/v4/stdlib"
)

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *application) connectToDB() (*sql.DB, error) {
	connection, err := openDB(app.DSN)
	if err != nil {
		return nil, err
	}

	log.Println("Connected to database successfully")

	return connection, nil
}

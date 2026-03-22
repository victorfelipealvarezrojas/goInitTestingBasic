package dbrepo

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

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

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbname   = "webapp_test"
	port     = "5435"
	dns      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"
)

var resource *dockertest.Resource
var poll *dockertest.Pool
var testdb *sql.DB

func TestMain(m *testing.M) {

	p, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	poll = p

	opt := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbname,
		},
		ExposedPorts: []string{port + ":5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	resource, err = poll.RunWithOptions(&opt)
	if err != nil {
		if resource != nil {
			_ = poll.Purge(resource)
		}
		log.Fatalf("Could not start resource: %s", err)
	}

	if err = poll.Retry(func() error {
		var err error
		testdb, err = sql.Open("pgx", fmt.Sprintf(dns, host, port, user, password, dbname))
		if err != nil {
			log.Printf("Could not connect to database: %s", err)
			return err
		}
		return testdb.Ping()
	}); err != nil {
		_ = poll.Purge(resource)
		log.Fatalf("Could not connect to database: %s", err)
	}

	// populate the database with test data
	err = createTables()
	if err != nil {
		log.Fatalf("Could not create tables: %s", err)
	}

	code := m.Run()

	// clean up
	if err = poll.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func createTables() error {
	tablesql, err := os.ReadFile("./testdata/users.sql")
	if err != nil {
		return err
	}

	_, err = testdb.Exec(string(tablesql))
	if err != nil {
		return err
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testdb.Ping()
	if err != nil {
		t.Errorf("Error pinging database: %s", err)
	}
}

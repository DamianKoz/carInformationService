package main

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Config struct {
	DB *sql.DB
}

func main() {

	db, err := connectToDB()
	if err != nil {
		panic(fmt.Sprintf("Could not connect db. Error: %v", err))
	}

	app := Config{
		DB: db,
	}

	err = app.initDB()
	if err != nil {
		panic(fmt.Sprintf("Could not connect db. Error: %v", err))
	}

	r := app.Routes()

	http.ListenAndServe(":8000", r)
}

func connectToDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost port=5432 user=admin password=postgres dbname=vehicles sslmode=disable timezone=UTC connect_timeout=5"
	}
	fmt.Printf("TESTING: %s", dsn)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (app *Config) initDB() error {
	err := app.createTables()
	if err != nil {
		return errors.New("Could not create tables in db. Error: " + err.Error())
	}

	return nil
}

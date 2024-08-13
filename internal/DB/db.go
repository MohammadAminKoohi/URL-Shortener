package DB

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
)

const (
	port     = 5432
	host     = "localhost"
	user     = "postgres"
	password = "Amin1385"
	dbname   = "urlshortener"
)

func DbInit() *sql.DB {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", user, password, host, port, dbname)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS urls (
    		id SERIAL PRIMARY KEY,
    		original_url TEXT NOT NULL,
    		shortened_url TEXT NOT NULL,
    		count INTEGER DEFAULT 0
		)
	`)
	if err != nil {
		log.Fatalf("Unable to create table: %v\n", err)
	} else {
		fmt.Println("Table is ready!")
	}

	return db
}

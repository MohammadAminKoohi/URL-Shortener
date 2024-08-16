package DB

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"log"
)

const (
	DB_PORT     = 5432
	DB_HOST     = "localhost"
	DB_USER     = "postgres"
	DB_PASSWORD = "Amineyk85"
	DB_NAME     = "urlshortener"
)
const (
	REDIS_PORT     = 6379
	REDIS_HOST     = "redis"
	REDIS_PASSWORD = ""
)

func DbInit() *sql.DB {
	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME)

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

func RedisInit() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", REDIS_HOST, REDIS_PORT),
		Password: REDIS_PASSWORD,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

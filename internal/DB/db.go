package DB

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
)

func DbInit() *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

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
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%v", redisHost, redisPort),
		Password: redisPassword,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	}
	return client
}

package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var Connection *pgxpool.Pool

// Connect function
func Initialize() (err error) {
	err = godotenv.Load()

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	Connection, err = pgxpool.New(context.Background(), connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	oniChan := make(chan bool, 1)
	go func(ch chan bool) {
		Connection.Ping(context.Background())
		ch <- true
	}(oniChan)

	select {
	case <-ctx.Done():
		fmt.Fprintln(os.Stderr, "Database Connection Timeout")
		os.Exit(1)
	case <-oniChan:
		log.Println("Database Connection Established!")
	}

	return nil
}

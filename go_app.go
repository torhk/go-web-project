package main

import (
	"os"
	"fmt"
	"log"
	"context"
	"net/http"
	"go_serv/db"
	"go_serv/routes"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbUrl := os.Getenv("DB_URL")
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
							dbUser, dbPass, dbUrl, dbName)
	pool, err := pgxpool.New(context.Background(), connStr)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Unable to create connection pool: %v\n", err)
        os.Exit(1)
    }
    defer pool.Close()

    user_db := &routes.Servers{
		Pg_db: db.New(pool),
	}
	err = pool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Database connection failed! Could not ping: %v", err)
	}
	
	fmt.Println("Ping successful! Database is online and reachable.")
	user_db.Pg_db.GetUserFromName(context.Background(),"thk")

	mux := routes.AddRoutes(user_db)
	log.Println("Starting server...")
	log.Fatal(http.ListenAndServe(":8090", mux))
}

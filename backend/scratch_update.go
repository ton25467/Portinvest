package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbURL := "postgresql://postgres:12345@localhost:5432/postgres"
	ctx := context.Background()

	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	hash := "$2a$10$3ElmzhAH05vEHtjx5ujcuettLgzzT/dq52ly3Lwx56sIwSi0qH26K"
	_, err = db.Exec(ctx, "UPDATE users SET password_hash = $1 WHERE email = 'user@example.com'", hash)
	if err != nil {
		fmt.Printf("Failed to update user password: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Password successfully updated in database!")
}

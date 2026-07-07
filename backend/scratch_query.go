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

	rows, err := db.Query(ctx, "SELECT id, email, password_hash FROM users")
	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		os.Exit(1)
	}
	defer rows.Close()

	found := false
	for rows.Next() {
		found = true
		var id, email, hash string
		err := rows.Scan(&id, &email, &hash)
		if err != nil {
			fmt.Printf("Row scan failed: %v\n", err)
			continue
		}
		fmt.Printf("User: %s | Email: %s | Hash: %s\n", id, email, hash)
	}

	if !found {
		fmt.Println("No users found in database.")
	}
}

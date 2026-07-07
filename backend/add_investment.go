package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	symbol := "K-JPX-A(A)"
	amount := 5000.0
	price := 31.9914
	quantity := amount / price

	var oldQty, oldAvgPrice float64
	err = db.QueryRow(ctx, "SELECT quantity, avg_buy_price FROM holdings WHERE symbol = $1 LIMIT 1", symbol).Scan(&oldQty, &oldAvgPrice)
	if err != nil {
		fmt.Printf("Failed to get holding: %v\n", err)
		os.Exit(1)
	}

	newQty := oldQty + quantity
	newCost := (oldQty * oldAvgPrice) + amount
	newAvgPrice := newCost / newQty

	_, err = db.Exec(ctx, "UPDATE holdings SET quantity = $1, avg_buy_price = $2 WHERE symbol = $3", newQty, newAvgPrice, symbol)
	if err != nil {
		fmt.Printf("Failed to update holding: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully updated %s! Old Qty: %.4f, New Qty: %.4f, Added Qty: %.4f, New Avg Price: %.4f\n", symbol, oldQty, newQty, quantity, newAvgPrice)
}

package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Holding struct {
	Symbol      string
	Name        string
	AssetType   string
	Quantity    float64
	AvgBuyPrice float64
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("No .env file found")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("DATABASE_URL is required")
		os.Exit(1)
	}

	ctx := context.Background()

	db, err := pgxpool.New(ctx, dbURL)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	var portfolioID string
	err = db.QueryRow(ctx, "SELECT id FROM portfolios LIMIT 1").Scan(&portfolioID)
	if err != nil {
		fmt.Printf("Failed to get portfolio: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Using portfolio: %s\n", portfolioID)

	// 1. Delete all old holdings
	_, err = db.Exec(ctx, "DELETE FROM holdings WHERE portfolio_id = $1", portfolioID)
	if err != nil {
		fmt.Printf("Failed to delete old holdings: %v\n", err)
		os.Exit(1)
	}

	// 2. Insert new exact holdings
	holdings := []Holding{
		{"KTB", "Krung Thai Bank", "stock", 200, 36.31},
		{"THAI", "Thai Airways International", "stock", 100, 6.21},
		{"TISCO", "Tisco Financial Group", "stock", 200, 118.45},
		{"SCBS&P500E", "SCB S&P 500 Index Fund", "etf", 231.2037, 43.2519},
		{"K-JPX-A(A)", "K Japan Share Index Fund", "etf", 16.1532, 30.9536},
		{"K-US500X-A(A)", "K US Equity Index Fund", "etf", 964.9653, 15.5446},
		{"K-USXNDQ-A(A)", "K US Nasdaq 100 Index Fund", "etf", 204.3673, 48.9315},
	}

	for _, h := range holdings {
		_, err := db.Exec(ctx, `
			INSERT INTO holdings (id, portfolio_id, symbol, name, asset_type, quantity, avg_buy_price, current_price)
			VALUES (gen_random_uuid(), $1, $2, $3, $4, $5, $6, $6)
		`, portfolioID, h.Symbol, h.Name, h.AssetType, h.Quantity, h.AvgBuyPrice)
		
		if err != nil {
			fmt.Printf("Failed to insert %s: %v\n", h.Symbol, err)
		} else {
			fmt.Printf("Inserted %s successfully.\n", h.Symbol)
		}
	}
}

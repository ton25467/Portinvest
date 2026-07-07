package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type FundItem struct {
	ID        string `json:"id"`
	ShortCode string `json:"short_code"`
}

func main() {
	url := "https://www.finnomena.com/fn3/api/fund/public/list"
	log.Printf("[Test] Fetching public list from: %s", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to fetch list: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read body: %v", err)
	}

	var list []FundItem
	if err := json.Unmarshal(body, &list); err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}

	log.Printf("Total funds found on Finnomena: %d", len(list))

	targets := []string{
		"k-jpx-a(a)", "k-us500x-a(a)", "k-usxndq-a(a)", "k-gold-a(a)", "k-wpbalanced", "scbs&p500e",
	}

	foundMap := make(map[string]string)
	for _, f := range list {
		code := strings.ToLower(f.ShortCode)
		for _, target := range targets {
			if code == target {
				foundMap[target] = f.ID
				log.Printf("FOUND FUND: %s -> ID: %s", f.ShortCode, f.ID)
			}
		}
	}

	outBytes, _ := json.MarshalIndent(foundMap, "", "  ")
	os.WriteFile("found_fund_ids.json", outBytes, 0644)
	log.Println("Test finished. Output saved to found_fund_ids.json")
}

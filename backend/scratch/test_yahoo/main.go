package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type YahooFinanceResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
				Currency           string  `json:"currency"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
}

func main() {
	stocks := []string{"TISCO.BK", "KTB.BK"}

	client := &http.Client{}

	for _, symbol := range stocks {
		url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", symbol)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")

		resp, err := client.Do(req)
		if err != nil {
			log.Printf("[Test] Failed to fetch %s: %v", symbol, err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[Test] Failed to read body %s: %v", symbol, err)
			continue
		}

		var yahooResp YahooFinanceResponse
		if err := json.Unmarshal(body, &yahooResp); err != nil {
			log.Printf("[Test] Failed to parse JSON for %s: %v", symbol, err)
			continue
		}

		if len(yahooResp.Chart.Result) > 0 {
			meta := yahooResp.Chart.Result[0].Meta
			log.Printf("STOCK %s -> Live Price: %.2f (%s)", symbol, meta.RegularMarketPrice, meta.Currency)
		} else {
			log.Printf("STOCK %s -> No quote details found", symbol)
		}
	}
}

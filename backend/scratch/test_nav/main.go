package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type NavResponse struct {
	Status bool `json:"status"`
	Data   struct {
		Navs []struct {
			Date  string  `json:"date"`
			Value float64 `json:"value"`
		} `json:"navs"`
	} `json:"data"`
}

func main() {
	funds := map[string]string{
		"K-JPX-A(A)":    "F00000X67S",
		"SCBS&P500E":    "F000013QNN",
		"K-USXNDQ-A(A)": "F0000143P4",
		"K-GOLD-A(A)":   "F000015I4E",
		"K-US500X-A(A)": "F00001C26R",
		"K-WPBALANCED":  "F00001CJHT",
	}

	for name, id := range funds {
		url := fmt.Sprintf("https://www.finnomena.com/fn3/api/fund/v2/public/funds/%s/nav/q?range=1M", id)
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("[Test] Failed to fetch %s: %v", name, err)
			continue
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("[Test] Failed to read body %s: %v", name, err)
			continue
		}

		var navResp NavResponse
		if err := json.Unmarshal(body, &navResp); err != nil {
			log.Printf("[Test] Failed to parse JSON for %s: %v", name, err)
			continue
		}

		if len(navResp.Data.Navs) > 0 {
			latest := navResp.Data.Navs[len(navResp.Data.Navs)-1]
			log.Printf("FUND %s -> Latest NAV: %.4f (Date: %s)", name, latest.Value, latest.Date)
		} else {
			log.Printf("FUND %s -> No NAV data found", name)
		}
	}
}

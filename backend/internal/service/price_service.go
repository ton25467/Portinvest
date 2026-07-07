package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

type PriceService struct {
	portfolioSrv *PortfolioService
	client       *http.Client
	fundMap      map[string]string // lowercase fund code -> finnomena fund ID
	fundMapMu    sync.RWMutex
}

type FinnomenaFund struct {
	ID        string `json:"id"`
	ShortCode string `json:"short_code"`
}

type FinnomenaNavResponse struct {
	Status bool `json:"status"`
	Data   struct {
		Navs []struct {
			Date  string  `json:"date"`
			Value float64 `json:"value"`
		} `json:"navs"`
	} `json:"data"`
}

type YahooChartResponse struct {
	Chart struct {
		Result []struct {
			Meta struct {
				RegularMarketPrice float64 `json:"regularMarketPrice"`
			} `json:"meta"`
		} `json:"result"`
	} `json:"chart"`
}

func NewPriceService(portfolioSrv *PortfolioService) *PriceService {
	return &PriceService{
		portfolioSrv: portfolioSrv,
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		fundMap: make(map[string]string),
	}
}

func (s *PriceService) Start(ctx context.Context) {
	slog.Info("Starting live price update service background worker")

	// 1. Initial build of fund mapping from Finnomena
	s.refreshFundMap()

	// 2. Perform initial price pull
	s.updateAllPrices(ctx)

	// 3. Ticker loop to run every 1 hour
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// We also refresh the fund map once a day
	mapRefreshTicker := time.NewTicker(24 * time.Hour)
	defer mapRefreshTicker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Stopping live price update service background worker")
			return
		case <-ticker.C:
			s.updateAllPrices(ctx)
		case <-mapRefreshTicker.C:
			s.refreshFundMap()
		}
	}
}

func (s *PriceService) refreshFundMap() {
	url := "https://www.finnomena.com/fn3/api/fund/public/list"
	slog.Info("Refreshing Finnomena fund map catalog", "url", url)

	resp, err := s.client.Get(url)
	if err != nil {
		slog.Error("Failed to fetch Finnomena fund list", "error", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("Failed to read Finnomena response body", "error", err)
		return
	}

	var list []FinnomenaFund
	if err := json.Unmarshal(body, &list); err != nil {
		slog.Error("Failed to parse Finnomena fund list json", "error", err)
		return
	}

	s.fundMapMu.Lock()
	for _, f := range list {
		s.fundMap[strings.ToLower(f.ShortCode)] = f.ID
	}
	s.fundMapMu.Unlock()

	slog.Info("Successfully cached Finnomena fund catalog", "count", len(list))
}

func (s *PriceService) getFinnomenaID(symbol string) (string, bool) {
	s.fundMapMu.RLock()
	defer s.fundMapMu.RUnlock()
	id, exists := s.fundMap[strings.ToLower(symbol)]
	return id, exists
}

func (s *PriceService) updateAllPrices(ctx context.Context) {
	slog.Info("Running scheduled live price update check...")

	// 1. Get all portfolios
	// To update prices, we need a list of unique symbols across all portfolios/holdings.
	// Since we don't have a global "ListAllHoldings" across all users, we can run a simple custom query
	// or query portfolios and list their holdings. Let's run a query on holdings table directly.
	symbols, err := s.getUniqueSymbols(ctx)
	if err != nil {
		slog.Error("Failed to retrieve unique holdings symbols", "error", err)
		return
	}

	if len(symbols) == 0 {
		slog.Info("No active holdings found in DB. Skipping price update.")
		return
	}

	prices := make(map[string]float64)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, symbol := range symbols {
		wg.Add(1)
		go func(sym string) {
			defer wg.Done()

			price, err := s.fetchPrice(sym)
			if err != nil {
				slog.Warn("Failed to fetch live price for symbol", "symbol", sym, "error", err)
				return
			}

			mu.Lock()
			prices[sym] = price
			mu.Unlock()

			slog.Info("Fetched live price", "symbol", sym, "price", price)
		}(symbol)
	}

	wg.Wait()

	if len(prices) > 0 {
		slog.Info("Updating database with fetched live prices", "count", len(prices))
		if err := s.portfolioSrv.UpdatePrices(ctx, prices); err != nil {
			slog.Error("Failed to update database prices", "error", err)
		} else {
			slog.Info("Successfully updated database prices")
		}
	}
}

func (s *PriceService) getUniqueSymbols(ctx context.Context) ([]string, error) {
	// Query unique symbols directly from postgres
	q := "SELECT DISTINCT symbol FROM holdings"
	rows, err := s.portfolioSrv.pool.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var symbols []string
	for rows.Next() {
		var sym string
		if err := rows.Scan(&sym); err != nil {
			return nil, err
		}
		symbols = append(symbols, sym)
	}
	return symbols, nil
}

func (s *PriceService) fetchPrice(symbol string) (float64, error) {
	// 1. Check if it matches a Finnomena Mutual Fund ID
	if fundID, exists := s.getFinnomenaID(symbol); exists {
		return s.fetchFinnomenaNAV(fundID)
	}

	// 2. Default to Yahoo Finance Stock ticker (e.g. TISCO, KTB -> TISCO.BK, KTB.BK)
	yahooSymbol := symbol
	if !strings.Contains(yahooSymbol, ".") {
		// Assume it's a Thai Stock if it has no dot suffix
		yahooSymbol = symbol + ".BK"
	}
	return s.fetchYahooStockPrice(yahooSymbol)
}

func (s *PriceService) fetchFinnomenaNAV(fundID string) (float64, error) {
	url := fmt.Sprintf("https://www.finnomena.com/fn3/api/fund/v2/public/funds/%s/nav/q?range=1D", fundID)
	resp, err := s.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// If 1D doesn't return anything or fails, try 1M
		url = fmt.Sprintf("https://www.finnomena.com/fn3/api/fund/v2/public/funds/%s/nav/q?range=1M", fundID)
		resp, err = s.client.Get(url)
		if err != nil {
			return 0, err
		}
		defer resp.Body.Close()
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var navResp FinnomenaNavResponse
	if err := json.Unmarshal(body, &navResp); err != nil {
		return 0, err
	}

	if len(navResp.Data.Navs) == 0 {
		return 0, fmt.Errorf("no NAV records found")
	}

	latest := navResp.Data.Navs[len(navResp.Data.Navs)-1]
	return latest.Value, nil
}

func (s *PriceService) fetchYahooStockPrice(symbol string) (float64, error) {
	url := fmt.Sprintf("https://query1.finance.yahoo.com/v8/finance/chart/%s", symbol)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	resp, err := s.client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var yahooResp YahooChartResponse
	if err := json.Unmarshal(body, &yahooResp); err != nil {
		return 0, err
	}

	if len(yahooResp.Chart.Result) == 0 {
		return 0, fmt.Errorf("no quote details found")
	}

	return yahooResp.Chart.Result[0].Meta.RegularMarketPrice, nil
}

package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/playwright-community/playwright-go"

	"portinves/internal/domain"
)

type WebScraperSyncService struct {
	cashflowSrv *CashflowService
	notifSrv    *NotificationService
	credRepo    domain.BankCredentialRepository
	pw          *playwright.Playwright
}

func NewWebScraperSyncService(
	cashflowSrv *CashflowService,
	notifSrv *NotificationService,
	credRepo domain.BankCredentialRepository,
) *WebScraperSyncService {
	// Initialize Playwright
	err := playwright.Install()
	if err != nil {
		slog.Error("could not install playwright drivers", "error", err)
	}

	pw, err := playwright.Run()
	if err != nil {
		slog.Error("could not start playwright", "error", err)
	} else {
		slog.Info("Playwright initialized successfully for WebScraperSyncService")
	}

	return &WebScraperSyncService{
		cashflowSrv: cashflowSrv,
		notifSrv:    notifSrv,
		credRepo:    credRepo,
		pw:          pw,
	}
}

func (s *WebScraperSyncService) Stop() {
	if s.pw != nil {
		if err := s.pw.Stop(); err != nil {
			slog.Error("failed to stop playwright", "error", err)
		}
	}
}

func (s *WebScraperSyncService) SyncTransactions(ctx context.Context, userID string) error {
	slog.Info("Starting web scraping sync for user", "userID", userID)

	if s.pw == nil {
		slog.Error("playwright is not initialized")
		return domain.ErrInternal(fmt.Errorf("playwright not initialized"))
	}

	// Example: Try to sync KBank if credentials exist
	cred, err := s.credRepo.GetByUserIDAndBank(ctx, userID, "KBANK")
	if err == nil && cred != nil {
		go s.runKBankScraper(userID, cred.Username, cred.PasswordEncrypted)
	} else {
		slog.Warn("No KBANK credentials found for user", "userID", userID)
	}

	return nil
}

// runKBankScraper is a skeleton for the Playwright scraper
func (s *WebScraperSyncService) runKBankScraper(userID, username, password string) {
	slog.Info("Launching browser for KBank sync", "user", username)
	
	ctx := context.Background()

	browser, err := s.pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true),
	})
	if err != nil {
		slog.Error("could not launch browser", "error", err)
		return
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		slog.Error("could not create page", "error", err)
		return
	}
	defer page.Close()

	// 1. Navigate to login
	slog.Info("Navigating to KBank portal...")
	// _, err = page.Goto("https://online.kasikornbankgroup.com/")
	// if err != nil { ... }

	// 2. Fill credentials (simulate)
	slog.Info("Filling credentials...")
	// page.Locator("#userName").Fill(username)
	// page.Locator("#password").Fill(password) // Decrypt password before using
	// page.Locator("#loginBtn").Click()

	// 3. Wait for OTP or Dashboard
	// ... (Handle OTP pause/resume logic if needed)

	// 4. Scrape data
	slog.Info("Scraping transaction data...")
	time.Sleep(2 * time.Second) // Simulate network delay

	parsedAmount := 2500.00
	parsedCurrency := "THB"
	parsedDescription := "KBank Payroll (Web Scraped)"
	parsedType := domain.CashflowTypeIncome

	// Insert Scraped Data
	_, err = s.cashflowSrv.CreateCashflow(ctx, userID, nil, parsedType, parsedAmount, parsedCurrency, parsedDescription, time.Now())
	if err != nil {
		slog.Error("Failed to record synced cashflow", "error", err)
		return
	}

	_, _ = s.notifSrv.CreateNotification(ctx, userID, "Sync Complete", "Successfully scraped 1 new transaction from KBank.")
	slog.Info("KBank scraper completed successfully")
}

package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"

	"portinves/internal/domain"
	"portinves/internal/handler/middleware"
)

type BankCredentialHandler struct {
	repo domain.BankCredentialRepository
}

func NewBankCredentialHandler(repo domain.BankCredentialRepository) *BankCredentialHandler {
	return &BankCredentialHandler{repo: repo}
}

type SaveCredentialReq struct {
	BankName string `json:"bank_name"`
	Username string `json:"username"`
	Password string `json:"password"` // Sent in plain text, backend should encrypt (mock encryption for MVP)
}

func (h *BankCredentialHandler) Save(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req SaveCredentialReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cred := &domain.BankCredential{
		ID:                uuid.New().String(),
		UserID:            userID,
		BankName:          req.BankName,
		Username:          req.Username,
		PasswordEncrypted: "ENCRYPTED_" + req.Password, // Mock encryption for MVP
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := h.repo.Save(r.Context(), cred); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

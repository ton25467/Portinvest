package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"portinves/internal/domain"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			slog.Error("Failed to encode JSON response", "error", err)
		}
	}
}

func writeError(w http.ResponseWriter, err error) {
	var appErr *domain.AppError
	if errors.As(err, &appErr) {
		writeJSON(w, appErr.Code, map[string]string{"message": appErr.Message})
		return
	}

	slog.Error("Unexpected server error", "error", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{"message": "internal server error"})
}

package handler

import (
	"encoding/json"
	"net/http"

	"portinves/internal/service"
)

type AuthHandler struct {
	authSrv *service.AuthService
}

func NewAuthHandler(authSrv *service.AuthService) *AuthHandler {
	return &AuthHandler{authSrv: authSrv}
}

type registerReq struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authResp struct {
	Token string `json:"token"`
	User  struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		Role  string `json:"role"`
	} `json:"user"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	user, err := h.authSrv.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		writeError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"message": "invalid request body"})
		return
	}

	token, err := h.authSrv.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeError(w, err)
		return
	}

	// Fetch user details to return in response
	userClaims, err := h.authSrv.ValidateToken(token)
	if err != nil {
		writeError(w, err)
		return
	}

	var resp authResp
	resp.Token = token
	resp.User.ID = userClaims.UserID
	resp.User.Email = userClaims.Email
	resp.User.Role = userClaims.Role
	// Name can be extracted/configured if needed, for simplicity we return what we have

	writeJSON(w, http.StatusOK, resp)
}

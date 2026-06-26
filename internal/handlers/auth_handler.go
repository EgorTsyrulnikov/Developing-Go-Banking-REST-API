package handlers

import (
	"encoding/json"
	"net/http"
	"bankapi/internal/config"
	"bankapi/internal/models"
	"bankapi/internal/services"
	"bankapi/pkg/logger"
)

type AuthHandler struct {
	Cfg *config.Config
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := services.RegisterUser(req)
	if err != nil {
		if err == services.ErrUserExists {
			http.Error(w, err.Error(), http.StatusConflict)
		} else {
			logger.Log.Errorf("Register error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	token, err := services.AuthenticateUser(req, h.Cfg)
	if err != nil {
		if err == services.ErrInvalidCredentials {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			logger.Log.Errorf("Login error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

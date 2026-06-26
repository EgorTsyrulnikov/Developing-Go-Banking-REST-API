package handlers

import (
	"encoding/json"
	"net/http"
	"bankapi/internal/config"
	"bankapi/internal/models"
	"bankapi/internal/services"
	"bankapi/pkg/logger"

	"github.com/gorilla/mux"
)

type CardHandler struct {
	Cfg *config.Config
}

func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req models.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	card, err := services.CreateCard(userID, req.AccountID, h.Cfg)
	if err != nil {
		if err == services.ErrUnauthorized || err == services.ErrAccountNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			logger.Log.Errorf("CreateCard error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(card)
}

func (h *CardHandler) GetCards(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	cards, err := services.GetCardsForAccount(userID, accountID)
	if err != nil {
		if err == services.ErrUnauthorized || err == services.ErrAccountNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			logger.Log.Errorf("GetCards error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cards)
}

package handlers

import (
	"encoding/json"
	"net/http"
	"bankapi/internal/models"
	"bankapi/internal/services"
	"bankapi/pkg/logger"
)

type AccountHandler struct{}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req models.CreateAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	account, err := services.CreateAccount(userID, req.Currency)
	if err != nil {
		logger.Log.Errorf("CreateAccount error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(account)
}

func (h *AccountHandler) Transfer(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req models.TransferRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := services.TransferFunds(userID, req.FromAccountID, req.ToAccountID, req.Amount)
	if err != nil {
		if err == services.ErrUnauthorized || err == services.ErrAccountNotFound || err == services.ErrInsufficientFunds {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			logger.Log.Errorf("Transfer error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	var req models.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := services.Deposit(req.AccountID, req.Amount)
	if err != nil {
		if err == services.ErrAccountNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			logger.Log.Errorf("Deposit error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

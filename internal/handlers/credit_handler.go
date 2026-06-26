package handlers

import (
	"encoding/json"
	"net/http"
	"bankapi/internal/models"
	"bankapi/internal/services"
	"bankapi/pkg/logger"

	"github.com/gorilla/mux"
)

type CreditHandler struct{}

func (h *CreditHandler) CreateCredit(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var req models.CreateCreditRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	credit, err := services.CreateCredit(userID, req.Amount, req.TermMonths)
	if err != nil {
		logger.Log.Errorf("CreateCredit error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(credit)
}

func (h *CreditHandler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	vars := mux.Vars(r)
	creditID := vars["creditId"]

	schedule, err := services.GetCreditSchedule(userID, creditID)
	if err != nil {
		if err == services.ErrUnauthorized {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			logger.Log.Errorf("GetSchedule error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(schedule)
}

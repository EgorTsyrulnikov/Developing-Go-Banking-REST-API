package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"bankapi/internal/services"
	"bankapi/pkg/logger"

	"github.com/gorilla/mux"
)

type AnalyticsHandler struct{}

func (h *AnalyticsHandler) GetAnalytics(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	vars := mux.Vars(r)
	accountID := vars["accountId"]

	analytics, err := services.GetAnalytics(userID, accountID)
	if err != nil {
		if err == services.ErrUnauthorized || err == services.ErrAccountNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			logger.Log.Errorf("GetAnalytics error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(analytics)
}

func (h *AnalyticsHandler) PredictBalance(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	vars := mux.Vars(r)
	accountID := vars["accountId"]
	daysStr := r.URL.Query().Get("days")
	
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 30 // default prediction 30 days
	}

	prediction, err := services.PredictBalance(userID, accountID, days)
	if err != nil {
		if err == services.ErrUnauthorized || err == services.ErrAccountNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			logger.Log.Errorf("PredictBalance error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(prediction)
}

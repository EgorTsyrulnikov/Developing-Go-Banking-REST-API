package routes

import (
	"bankapi/internal/config"
	"bankapi/internal/handlers"
	"bankapi/internal/middleware"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router, cfg *config.Config) {
	authHandler := &handlers.AuthHandler{Cfg: cfg}
	accountHandler := &handlers.AccountHandler{}
	cardHandler := &handlers.CardHandler{Cfg: cfg}
	creditHandler := &handlers.CreditHandler{}
	analyticsHandler := &handlers.AnalyticsHandler{}

	// Public routes
	r.HandleFunc("/register", authHandler.Register).Methods("POST")
	r.HandleFunc("/login", authHandler.Login).Methods("POST")

	// Protected routes
	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(middleware.AuthMiddleware(cfg))

	// Accounts
	authRouter.HandleFunc("/accounts", accountHandler.CreateAccount).Methods("POST")
	authRouter.HandleFunc("/accounts/deposit", accountHandler.Deposit).Methods("POST")
	authRouter.HandleFunc("/transfer", accountHandler.Transfer).Methods("POST")

	// Cards
	authRouter.HandleFunc("/cards", cardHandler.CreateCard).Methods("POST")
	authRouter.HandleFunc("/accounts/{accountId}/cards", cardHandler.GetCards).Methods("GET")

	// Credits
	authRouter.HandleFunc("/credits", creditHandler.CreateCredit).Methods("POST")
	authRouter.HandleFunc("/credits/{creditId}/schedule", creditHandler.GetSchedule).Methods("GET")

	// Analytics
	authRouter.HandleFunc("/accounts/{accountId}/analytics", analyticsHandler.GetAnalytics).Methods("GET")
	authRouter.HandleFunc("/accounts/{accountId}/predict", analyticsHandler.PredictBalance).Methods("GET")
}

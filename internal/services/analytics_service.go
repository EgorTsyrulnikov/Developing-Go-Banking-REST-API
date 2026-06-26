package services

import (
	"bankapi/internal/models"
	"bankapi/internal/repositories"
)

func GetAnalytics(userID, accountID string) (*models.AnalyticsResponse, error) {
	acc, err := repositories.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, ErrAccountNotFound
	}
	if acc.UserID != userID {
		return nil, ErrUnauthorized
	}

	txs, err := repositories.GetTransactionsByAccount(accountID)
	if err != nil {
		return nil, err
	}

	var income, expense float64
	for _, tx := range txs {
		if tx.ToAccountID != nil && *tx.ToAccountID == accountID {
			income += tx.Amount
		} else if tx.FromAccountID != nil && *tx.FromAccountID == accountID {
			expense += tx.Amount
		}
	}

	totalCredit, err := repositories.GetUnpaidCreditLoadByUserID(userID)
	if err != nil {
		return nil, err
	}

	return &models.AnalyticsResponse{
		TotalIncome:  income,
		TotalExpense: expense,
		TotalCredit:  totalCredit,
	}, nil
}

func PredictBalance(userID, accountID string, days int) (*models.PredictBalanceResponse, error) {
	if days > 365 {
		days = 365
	}
	acc, err := repositories.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if acc == nil {
		return nil, ErrAccountNotFound
	}
	if acc.UserID != userID {
		return nil, ErrUnauthorized
	}

	upcomingCreditPayments, err := repositories.GetUpcomingPaymentsByUserID(userID, days)
	if err != nil {
		return nil, err
	}

	predicted := acc.Balance - upcomingCreditPayments

	return &models.PredictBalanceResponse{
		PredictedBalance: predicted,
	}, nil
}

package models

type AnalyticsResponse struct {
	TotalIncome  float64 `json:"total_income"`
	TotalExpense float64 `json:"total_expense"`
	TotalCredit  float64 `json:"total_credit_debt"`
}

type PredictBalanceResponse struct {
	PredictedBalance float64 `json:"predicted_balance"`
}

package models

import "time"

type Credit struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Amount         float64   `json:"amount"`
	InterestRate   float64   `json:"interest_rate"`
	TermMonths     int       `json:"term_months"`
	MonthlyPayment float64   `json:"monthly_payment"`
	CreatedAt      time.Time `json:"created_at"`
}

type PaymentSchedule struct {
	ID          string    `json:"id"`
	CreditID    string    `json:"credit_id"`
	PaymentDate time.Time `json:"payment_date"`
	Amount      float64   `json:"amount"`
	IsPaid      bool      `json:"is_paid"`
	Penalty     float64   `json:"penalty"`
}

type CreateCreditRequest struct {
	Amount     float64 `json:"amount"`
	TermMonths int     `json:"term_months"`
}

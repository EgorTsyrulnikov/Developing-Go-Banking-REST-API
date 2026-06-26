package services

import (
	"errors"
	"math"
	"time"
	"bankapi/internal/models"
	"bankapi/internal/repositories"
	"bankapi/pkg/cbr"
)

func CreateCredit(userID string, amount float64, termMonths int) (*models.Credit, error) {
	rate, err := cbr.GetCentralBankRate()
	if err != nil {
		// Fallback to a default rate if CBR is unavailable
		rate = 15.0
	}

	monthlyRate := rate / 100 / 12
	annuityFactor := (monthlyRate * math.Pow(1+monthlyRate, float64(termMonths))) / (math.Pow(1+monthlyRate, float64(termMonths)) - 1)
	monthlyPayment := amount * annuityFactor

	credit := &models.Credit{
		UserID:         userID,
		Amount:         amount,
		InterestRate:   rate,
		TermMonths:     termMonths,
		MonthlyPayment: math.Round(monthlyPayment*100) / 100, // Round to 2 decimals
	}

	tx, err := repositories.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	err = repositories.CreateCreditTx(tx, credit)
	if err != nil {
		return nil, err
	}

	// Generate payment schedules
	for i := 1; i <= termMonths; i++ {
		schedule := &models.PaymentSchedule{
			CreditID:    credit.ID,
			PaymentDate: time.Now().AddDate(0, i, 0),
			Amount:      credit.MonthlyPayment,
			IsPaid:      false,
		}
		err = repositories.CreatePaymentScheduleTx(tx, schedule)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return credit, nil
}

func GetCreditSchedule(userID, creditID string) ([]models.PaymentSchedule, error) {
	credit, err := repositories.GetCreditByID(creditID)
	if err != nil {
		return nil, err
	}
	if credit == nil {
		return nil, errors.New("credit not found")
	}
	if credit.UserID != userID {
		return nil, ErrUnauthorized
	}

	return repositories.GetPaymentSchedulesByCreditID(creditID)
}

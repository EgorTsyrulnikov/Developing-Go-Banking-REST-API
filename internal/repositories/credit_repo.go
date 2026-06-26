package repositories

import (
	"database/sql"
	"bankapi/internal/models"
	"time"
)

func CreateCreditTx(tx *sql.Tx, credit *models.Credit) error {
	query := `INSERT INTO credits (user_id, amount, interest_rate, term_months, monthly_payment) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return tx.QueryRow(query, credit.UserID, credit.Amount, credit.InterestRate, credit.TermMonths, credit.MonthlyPayment).Scan(&credit.ID, &credit.CreatedAt)
}

func CreatePaymentScheduleTx(tx *sql.Tx, schedule *models.PaymentSchedule) error {
	query := `INSERT INTO payment_schedules (credit_id, payment_date, amount, is_paid) 
	VALUES ($1, $2, $3, $4) RETURNING id`
	return tx.QueryRow(query, schedule.CreditID, schedule.PaymentDate, schedule.Amount, schedule.IsPaid).Scan(&schedule.ID)
}

func GetCreditByID(id string) (*models.Credit, error) {
	credit := &models.Credit{}
	query := `SELECT id, user_id, amount, interest_rate, term_months, monthly_payment, created_at FROM credits WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&credit.ID, &credit.UserID, &credit.Amount, &credit.InterestRate, &credit.TermMonths, &credit.MonthlyPayment, &credit.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return credit, nil
}

func GetPaymentSchedulesByCreditID(creditID string) ([]models.PaymentSchedule, error) {
	query := `SELECT id, credit_id, payment_date, amount, is_paid, penalty FROM payment_schedules WHERE credit_id = $1 ORDER BY payment_date ASC`
	rows, err := DB.Query(query, creditID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PaymentSchedule
	for rows.Next() {
		var s models.PaymentSchedule
		if err := rows.Scan(&s.ID, &s.CreditID, &s.PaymentDate, &s.Amount, &s.IsPaid, &s.Penalty); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func GetOverduePayments(cutoff time.Time) ([]models.PaymentSchedule, error) {
	query := `SELECT id, credit_id, payment_date, amount, is_paid, penalty FROM payment_schedules 
	WHERE is_paid = FALSE AND payment_date <= $1`
	rows, err := DB.Query(query, cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []models.PaymentSchedule
	for rows.Next() {
		var s models.PaymentSchedule
		if err := rows.Scan(&s.ID, &s.CreditID, &s.PaymentDate, &s.Amount, &s.IsPaid, &s.Penalty); err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func UpdatePaymentSchedulePenalty(scheduleID string, penalty float64) error {
	query := `UPDATE payment_schedules SET penalty = $1 WHERE id = $2`
	_, err := DB.Exec(query, penalty, scheduleID)
	return err
}

func MarkPaymentPaid(scheduleID string) error {
	query := `UPDATE payment_schedules SET is_paid = TRUE WHERE id = $1`
	_, err := DB.Exec(query, scheduleID)
	return err
}

func GetUnpaidCreditLoadByUserID(userID string) (float64, error) {
	query := `SELECT COALESCE(SUM(ps.amount + ps.penalty), 0)
	          FROM payment_schedules ps
	          JOIN credits c ON ps.credit_id = c.id
	          WHERE c.user_id = $1 AND ps.is_paid = FALSE`
	var total float64
	err := DB.QueryRow(query, userID).Scan(&total)
	return total, err
}

func GetUpcomingPaymentsByUserID(userID string, days int) (float64, error) {
	cutoff := time.Now().AddDate(0, 0, days)
	query := `SELECT COALESCE(SUM(ps.amount + ps.penalty), 0)
	          FROM payment_schedules ps
	          JOIN credits c ON ps.credit_id = c.id
	          WHERE c.user_id = $1 AND ps.is_paid = FALSE AND ps.payment_date <= $2`
	var total float64
	err := DB.QueryRow(query, userID, cutoff).Scan(&total)
	return total, err
}

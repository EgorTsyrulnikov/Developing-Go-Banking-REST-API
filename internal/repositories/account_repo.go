package repositories

import (
	"database/sql"
	"bankapi/internal/models"
)

func CreateAccount(account *models.Account) error {
	query := `INSERT INTO accounts (user_id, currency) VALUES ($1, $2) RETURNING id, balance, created_at`
	return DB.QueryRow(query, account.UserID, account.Currency).Scan(&account.ID, &account.Balance, &account.CreatedAt)
}

func GetAccountByID(id string) (*models.Account, error) {
	account := &models.Account{}
	query := `SELECT id, user_id, balance, currency, created_at FROM accounts WHERE id = $1`
	err := DB.QueryRow(query, id).Scan(&account.ID, &account.UserID, &account.Balance, &account.Currency, &account.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return account, nil
}

func UpdateAccountBalanceTx(tx *sql.Tx, accountID string, amount float64) error {
	query := `UPDATE accounts SET balance = balance + $1 WHERE id = $2`
	_, err := tx.Exec(query, amount, accountID)
	return err
}

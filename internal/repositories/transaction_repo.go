package repositories

import (
	"database/sql"
	"bankapi/internal/models"
)

func CreateTransactionTx(tx *sql.Tx, transaction *models.Transaction) error {
	query := `INSERT INTO transactions (from_account_id, to_account_id, amount, type) VALUES ($1, $2, $3, $4) RETURNING id, created_at`
	return tx.QueryRow(query, transaction.FromAccountID, transaction.ToAccountID, transaction.Amount, transaction.Type).Scan(&transaction.ID, &transaction.CreatedAt)
}

func GetTransactionsByAccount(accountID string) ([]models.Transaction, error) {
	query := `SELECT id, from_account_id, to_account_id, amount, type, created_at 
	FROM transactions WHERE from_account_id = $1 OR to_account_id = $1 ORDER BY created_at DESC`
	rows, err := DB.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txs []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.FromAccountID, &t.ToAccountID, &t.Amount, &t.Type, &t.CreatedAt); err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, nil
}

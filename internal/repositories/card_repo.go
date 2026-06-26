package repositories

import (
	"bankapi/internal/models"
)

func CreateCard(card *models.Card) error {
	query := `INSERT INTO cards (account_id, card_number_encrypted, card_number_hmac, expiration_date, cvv_hash) 
	VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	return DB.QueryRow(query, card.AccountID, card.CardNumberEncrypted, card.CardNumberHMAC, card.ExpirationDate, card.CVVHash).Scan(&card.ID, &card.CreatedAt)
}

func GetCardsByAccountID(accountID string) ([]models.Card, error) {
	query := `SELECT id, account_id, card_number_encrypted, card_number_hmac, expiration_date, cvv_hash, created_at FROM cards WHERE account_id = $1`
	rows, err := DB.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cards []models.Card
	for rows.Next() {
		var c models.Card
		if err := rows.Scan(&c.ID, &c.AccountID, &c.CardNumberEncrypted, &c.CardNumberHMAC, &c.ExpirationDate, &c.CVVHash, &c.CreatedAt); err != nil {
			return nil, err
		}
		cards = append(cards, c)
	}
	return cards, nil
}

package models

import "time"

type Card struct {
	ID                  string    `json:"id"`
	AccountID           string    `json:"account_id"`
	CardNumberEncrypted string    `json:"-"`
	CardNumberDecrypted string    `json:"card_number,omitempty"`
	CardNumberHMAC      string    `json:"-"`
	ExpirationDate      string    `json:"expiration_date"`
	CVVHash             string    `json:"-"`
	CreatedAt           time.Time `json:"created_at"`
}

type CreateCardRequest struct {
	AccountID string `json:"account_id"`
}

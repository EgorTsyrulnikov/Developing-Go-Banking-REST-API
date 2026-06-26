package services

import (
	"errors"
	"bankapi/internal/models"
	"bankapi/internal/repositories"
)

var ErrInsufficientFunds = errors.New("insufficient funds")
var ErrAccountNotFound = errors.New("account not found")
var ErrUnauthorized = errors.New("unauthorized access to account")

func CreateAccount(userID, currency string) (*models.Account, error) {
	if currency == "" {
		currency = "RUB"
	}
	account := &models.Account{
		UserID:   userID,
		Currency: currency,
	}
	err := repositories.CreateAccount(account)
	if err != nil {
		return nil, err
	}
	return account, nil
}

func GetAccountByID(accountID string) (*models.Account, error) {
	account, err := repositories.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	return account, nil
}

func TransferFunds(userID, fromAccountID, toAccountID string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	fromAcc, err := repositories.GetAccountByID(fromAccountID)
	if err != nil {
		return err
	}
	if fromAcc == nil {
		return ErrAccountNotFound
	}
	if fromAcc.UserID != userID {
		return ErrUnauthorized
	}

	toAcc, err := repositories.GetAccountByID(toAccountID)
	if err != nil {
		return err
	}
	if toAcc == nil {
		return ErrAccountNotFound
	}

	if fromAcc.Balance < amount {
		return ErrInsufficientFunds
	}

	tx, err := repositories.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = repositories.UpdateAccountBalanceTx(tx, fromAccountID, -amount)
	if err != nil {
		return err
	}

	err = repositories.UpdateAccountBalanceTx(tx, toAccountID, amount)
	if err != nil {
		return err
	}

	transaction := &models.Transaction{
		FromAccountID: &fromAccountID,
		ToAccountID:   &toAccountID,
		Amount:        amount,
		Type:          "transfer",
	}

	err = repositories.CreateTransactionTx(tx, transaction)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func Deposit(accountID string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}

	acc, err := repositories.GetAccountByID(accountID)
	if err != nil {
		return err
	}
	if acc == nil {
		return ErrAccountNotFound
	}

	tx, err := repositories.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = repositories.UpdateAccountBalanceTx(tx, accountID, amount)
	if err != nil {
		return err
	}

	transaction := &models.Transaction{
		ToAccountID: &accountID,
		Amount:      amount,
		Type:        "deposit",
	}

	err = repositories.CreateTransactionTx(tx, transaction)
	if err != nil {
		return err
	}

	return tx.Commit()
}

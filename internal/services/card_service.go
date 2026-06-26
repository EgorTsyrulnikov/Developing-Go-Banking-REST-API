package services

import (
	"fmt"
	"math/rand"
	"time"
	"bankapi/internal/config"
	"bankapi/internal/models"
	"bankapi/internal/repositories"
	"bankapi/pkg/crypto"
)

func generateLuhnCardNumber() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	prefix := "400000" // Visa like
	for i := 0; i < 9; i++ {
		prefix += fmt.Sprintf("%d", r.Intn(10))
	}

	sum := 0
	for i := 0; i < len(prefix); i++ {
		digit := int(prefix[i] - '0')
		if i%2 == 0 {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	checkDigit := (10 - (sum % 10)) % 10
	return fmt.Sprintf("%s%d", prefix, checkDigit)
}

func CreateCard(userID, accountID string, cfg *config.Config) (*models.Card, error) {
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

	cardNumber := generateLuhnCardNumber()
	expirationDate := time.Now().AddDate(3, 0, 0).Format("01/06") // MM/YY

	// Generate CVV 
	cvv := fmt.Sprintf("%03d", rand.Intn(1000))

	encryptedNum, err := crypto.EncryptPGP(cardNumber)
	if err != nil {
		return nil, err
	}

	hmacNum := crypto.ComputeHMAC(cardNumber, cfg.JWTSecret)
	
	cvvHash, err := crypto.HashCVV(cvv)
	if err != nil {
		return nil, err
	}

	card := &models.Card{
		AccountID:           accountID,
		CardNumberEncrypted: encryptedNum,
		CardNumberHMAC:      hmacNum,
		ExpirationDate:      expirationDate,
		CVVHash:             cvvHash,
	}

	err = repositories.CreateCard(card)
	if err != nil {
		return nil, err
	}

	// Attach decrypted value only for the initial response
	card.CardNumberDecrypted = cardNumber
	// In a real app we might return CVV here just once, but we will omit it for simplicity

	return card, nil
}

func GetCardsForAccount(userID, accountID string) ([]models.Card, error) {
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

	cards, err := repositories.GetCardsByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	for i := range cards {
		decrypted, err := crypto.DecryptPGP(cards[i].CardNumberEncrypted)
		if err == nil {
			cards[i].CardNumberDecrypted = decrypted
		}
	}

	return cards, nil
}

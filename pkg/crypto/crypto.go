package crypto

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"

	"github.com/ProtonMail/go-crypto/openpgp"
	"golang.org/x/crypto/bcrypt"
)

var (
	Entity *openpgp.Entity
)

func InitPGP() error {
	// Generate a new key for testing/dev purposes
	var err error
	Entity, err = openpgp.NewEntity("Bank", "Test", "bank@example.com", nil)
	return err
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashCVV(cvv string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(cvv), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckCVVHash(cvv, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(cvv))
	return err == nil
}

func ComputeHMAC(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func EncryptPGP(data string) (string, error) {
	if Entity == nil {
		return "", fmt.Errorf("PGP entity not initialized")
	}

	buf := new(bytes.Buffer)
	w, err := openpgp.Encrypt(buf, []*openpgp.Entity{Entity}, nil, nil, nil)
	if err != nil {
		return "", err
	}
	_, err = w.Write([]byte(data))
	if err != nil {
		return "", err
	}
	w.Close()

	return hex.EncodeToString(buf.Bytes()), nil
}

func DecryptPGP(encryptedHex string) (string, error) {
	if Entity == nil {
		return "", fmt.Errorf("PGP entity not initialized")
	}

	encryptedData, err := hex.DecodeString(encryptedHex)
	if err != nil {
		return "", err
	}

	buf := bytes.NewReader(encryptedData)
	entityList := openpgp.EntityList{Entity}
	md, err := openpgp.ReadMessage(buf, entityList, nil, nil)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := io.ReadAll(md.UnverifiedBody)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}

package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"math/big"
)

// Generate a random verification token with the specified length
func GenerateVerificationToken(length uint) (string, error) {
	randomBytes := make([]byte, (length+1)/2)

	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	verification_token := hex.EncodeToString(randomBytes)[:length]
	return verification_token, nil
}

// Generate password for students
func GeneratePassword(length int) (string, error) {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+,.?/:;{}[]"
	charsetLength := big.NewInt(int64(len(charset)))

	randomString := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		randomString[i] = charset[randomIndex.Int64()]
	}
	return string(randomString), nil
}

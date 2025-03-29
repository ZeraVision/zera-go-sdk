package wallet

import (
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/tyler-smith/go-bip39"
)

// GenerateMnemonic generates a BIP39 mnemonic phrase based on the given entropy strength.
func GenerateMnemonic(strength int) (string, error) {
	if strength%32 != 0 || strength < 128 || strength > 256 {
		return "", errors.New("entropy strength must be one of 128, 160, 192, 224, or 256 bits")
	}
	entropy, err := bip39.NewEntropy(strength)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+[]{}|;:,.<>?/"

// GenerateRandomString generates a random string based on the given length and character set.
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be greater than 0")
	}

	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := range result {
		randomInt, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = charset[randomInt.Int64()]
	}

	return string(result), nil
}

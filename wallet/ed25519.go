package wallet

import (
	"crypto/sha256"
	"errors"

	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	ed25519 "github.com/teserakt-io/golang-ed25519"
)

// GenerateKeyPairLibsodium-compatible: deterministically creates an Ed25519 keypair from a 32-byte seed.
func GenerateKeyPairLibsodium(seed []byte) ([]byte, []byte, error) {
	if len(seed) != 32 {
		return nil, nil, errors.New("seed must be exactly 32 bytes")
	}

	privateKey := ed25519.NewKeyFromSeed(seed)
	publicKey := privateKey.Public().(ed25519.PublicKey)

	return privateKey, publicKey, nil
}

// GenerateEd25519 generates an Ed25519 keypair and returns the encoded values with optional mnemonic seed.
func GenerateEd25519(mnemonic string, hashAlg helper.HashType, keyType helper.KeyType) (string, string, string, error) {
	if mnemonic == "" {
		var err error
		mnemonic, err = GenerateRandomString(1000)
		if err != nil {
			return "", "", "", errors.New("failed to generate random entropy")
		}
	}

	// Replace this with real Argon2 or other strong hashing if needed
	seed := sha256.Sum256([]byte(mnemonic))

	privateKey, rawPublicKey, err := GenerateKeyPairLibsodium(seed[:])
	if err != nil {
		return "", "", "", err
	}

	publicKey, b58Address, err := GetWalletAddress(rawPublicKey, hashAlg, keyType)
	if err != nil {
		return "", "", "", err
	}

	b58PublicKey := transcode.Base58Encode(rawPublicKey)

	underscoreCount := 0
	secondUnderscoreIndex := -1
	for i, b := range publicKey {
		if b == '_' {
			underscoreCount++
			if underscoreCount == 2 {
				secondUnderscoreIndex = i
				break
			}
		}
	}

	if secondUnderscoreIndex != -1 {
		prefix := string(publicKey[:secondUnderscoreIndex+1])
		b58PublicKey = prefix + b58PublicKey
	}

	b58Private := transcode.Base58Encode(privateKey)

	return b58Private, b58PublicKey, b58Address, nil
}

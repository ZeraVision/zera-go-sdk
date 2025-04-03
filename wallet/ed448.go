package wallet

import (
	"errors"
	"fmt"

	generichash "github.com/GoKillers/libsodium-go/cryptogenerichash"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zn-wallet-manager/transcode"
	"github.com/cloudflare/circl/sign/ed448"
)

// GenerateKeyPairEd448 generates an Ed448 key pair using circl.
func GenerateKeyPairEd448(seed []byte) ([]byte, []byte, error) {
	if len(seed) != ed448.SeedSize {
		return nil, nil, fmt.Errorf("seed must be exactly %d bytes", ed448.SeedSize)
	}

	// Generate private and public key
	privateKey := ed448.NewKeyFromSeed(seed)
	publicKey := privateKey.Public().(ed448.PublicKey)

	return privateKey.Seed(), publicKey[:], nil
}

// GenerateEd448 generates an Ed448 key pair and hashes the public key with the specified algorithm.
func GenerateEd448(mnemonic string, hashAlg helper.HashType, keyType helper.KeyType) (string, string, string, error) {
	if len(mnemonic) < 12 {
		var err error
		mnemonic, err = GenerateRandomString(1000)

		if err != nil {
			return "", "", "", errors.New("failed to generate random entropy")
		}
	}

	seed, retCode := generichash.CryptoGenericHash(ed448.SeedSize, []byte(mnemonic), nil)

	if retCode != 0 {
		return "", "", "", errors.New("libsodium: failed to generate seed")
	}

	privateKey, rawPublicKey, err := GenerateKeyPairEd448(seed)
	if err != nil {
		return "", "", "", err
	}

	publicKey, b58Address, err := GetWalletAddress(rawPublicKey, hashAlg, keyType)
	if err != nil {
		return "", "", "", err
	}

	b58PublicKey := transcode.Base58Encode(rawPublicKey)

	// Find the index of the second underscore
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

	// Prepend everything up to and including the second underscore to the b58PublicKey
	if secondUnderscoreIndex != -1 {
		prefix := string(publicKey[:secondUnderscoreIndex+1])
		b58PublicKey = prefix + b58PublicKey
	}

	b58Private := transcode.Base58Encode(privateKey)

	return b58Private, b58PublicKey, b58Address, nil
}

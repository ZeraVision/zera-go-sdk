package wallet

import (
	"errors"

	generichash "github.com/GoKillers/libsodium-go/cryptogenerichash"
	"github.com/GoKillers/libsodium-go/cryptosign"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/transcode"
)

// GenerateKeyPairLibsodium uses libsodium to generate an Ed25519 key pair from a seed.
func GenerateKeyPairLibsodium(seed []byte) ([]byte, []byte, error) {
	if len(seed) != 32 {
		return nil, nil, errors.New("seed must be exactly 32 bytes")
	}

	// Generate key pair using libsodium
	privateKey, publicKey, ret := cryptosign.CryptoSignSeedKeyPair(seed)
	if ret != 0 {
		return nil, nil, errors.New("libsodium: failed to generate key pair")
	}

	return privateKey, publicKey, nil
}

// GenerateEd25519 generates an Ed25519 key pair and hashes the public key with the specified algorithm.
func GenerateEd25519(mnemonic string, hashAlg helper.HashType, keyType helper.KeyType) (string, string, string, error) {

	// If empty, generate random entropy not based on BIP39
	if mnemonic == "" {
		var err error
		mnemonic, err = GenerateRandomString(1000)

		if err != nil {
			return "", "", "", errors.New("failed to generate random entropy")
		}
	}

	seed, retCode := generichash.CryptoGenericHash(32, []byte(mnemonic), nil)

	if retCode != 0 {
		return "", "", "", errors.New("libsodium: failed to generate seed")
	}

	privateKey, rawPublicKey, err := GenerateKeyPairLibsodium(seed)
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

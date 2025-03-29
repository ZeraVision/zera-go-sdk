package sign

import (
	"errors"
	"fmt"

	"github.com/GoKillers/libsodium-go/cryptosign"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zn-wallet-manager/transcode"
	"github.com/cloudflare/circl/sign/ed448"
)

// SignTransaction signs a transaction payload using the Ed25519 private key.
func Sign(privateKeyBase58 string, payload []byte, keyType helper.KeyType) ([]byte, error) {
	if len(payload) == 0 {
		return nil, errors.New("payload cannot be empty")
	}

	// Decode the Base58-encoded private key
	privateKey, err := transcode.Base58Decode(privateKeyBase58)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	var signature []byte

	if keyType == helper.ED25519 {
		// Ensure the private key length matches Ed25519 requirements
		if len(privateKey) != cryptosign.CryptoSignSecretKeyBytes() {
			return nil, errors.New("invalid private key length")
		}

		var retCode int
		signature, retCode = cryptosign.CryptoSignDetached(payload, privateKey)
		if retCode != 0 {
			return nil, errors.New("libsodium: failed to sign transaction")
		}
	} else if keyType == helper.ED448 {
		// Ensure the private key length matches Ed448 requirements
		if len(privateKey) != 57 {
			return nil, errors.New("invalid private key length")
		}

		privKey := ed448.NewKeyFromSeed(privateKey)
		signature := ed448.Sign(privKey, payload, "")
		return signature, nil
	} else {
		return nil, errors.New("unsupported key type")
	}

	return signature, nil
}

// Verify verifies the signature of a transaction payload using the public key and the key type.
func Verify(publicKeyBase58 string, payload []byte, signature []byte) (bool, error) {
	_, publicKeyByte, _, err := transcode.Base58DecodePublicKey(publicKeyBase58)

	if err != nil {
		return false, fmt.Errorf("could not decode public key: %v", err)
	}

	var keyType helper.KeyType
	if len(publicKeyBase58) > 0 {
		if publicKeyBase58[0] == 'A' {
			keyType = helper.ED25519
		} else if publicKeyBase58[0] == 'B' {
			keyType = helper.ED448
		} else {
			return false, errors.New("unsupported key type")
		}
	} else {
		return false, errors.New("public key is empty")
	}

	if len(payload) == 0 {
		return false, errors.New("payload cannot be empty")
	}

	if len(signature) == 0 {
		return false, errors.New("signature cannot be empty")
	}

	if keyType == helper.ED25519 {
		// Ensure the public key length matches Ed25519 requirements
		if len(publicKeyByte) != cryptosign.CryptoSignPublicKeyBytes() {
			return false, errors.New("invalid public key length for ED25519")
		}

		// Verify the signature
		verified := cryptosign.CryptoSignVerifyDetached(signature, payload, publicKeyByte)
		if verified != 0 {
			return false, errors.New("libsodium: signature verification failed")
		}

		return true, nil
	} else if keyType == helper.ED448 {
		// Ensure the public key length matches Ed448 requirements
		if len(publicKeyByte) != ed448.PublicKeySize {
			return false, errors.New("invalid public key length for ED448")
		}

		verified := ed448.Verify(publicKeyByte, payload, signature, "")
		if !verified {
			return false, errors.New("ED448: signature verification failed")
		}

		return true, nil
	} else {
		return false, errors.New("unsupported key type")
	}
}

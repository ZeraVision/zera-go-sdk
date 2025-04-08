package sign

import (
	"errors"
	"fmt"

	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zn-wallet-manager/transcode"
	"github.com/cloudflare/circl/sign/ed448"
	ed25519 "github.com/teserakt-io/golang-ed25519"
)

// SignTransaction signs a transaction payload using Ed25519 or Ed448.
func Sign(privateKeyBase58 string, payload []byte, keyType helper.KeyType) ([]byte, error) {
	if len(payload) == 0 {
		return nil, errors.New("payload cannot be empty")
	}

	privateKey, err := transcode.Base58Decode(privateKeyBase58)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	switch keyType {
	case helper.ED25519:
		if len(privateKey) != ed25519.PrivateKeySize {
			return nil, errors.New("invalid private key length for ED25519")
		}
		signature := ed25519.Sign(privateKey, payload)
		return signature, nil

	case helper.ED448:
		if len(privateKey) != 57 {
			return nil, errors.New("invalid private key length for ED448")
		}
		privKey := ed448.NewKeyFromSeed(privateKey)
		signature := ed448.Sign(privKey, payload, "")
		return signature, nil

	default:
		return nil, errors.New("unsupported key type")
	}
}

// Verify checks the signature of a payload using the given public key.
func Verify(publicKeyBase58 string, payload []byte, signature []byte) (bool, error) {
	_, publicKeyByte, _, err := transcode.Base58DecodePublicKey(publicKeyBase58)
	if err != nil {
		return false, fmt.Errorf("could not decode public key: %v", err)
	}

	if len(payload) == 0 {
		return false, errors.New("payload cannot be empty")
	}

	if len(signature) == 0 {
		return false, errors.New("signature cannot be empty")
	}

	var keyType helper.KeyType
	if len(publicKeyBase58) > 0 {
		switch publicKeyBase58[0] {
		case 'A':
			keyType = helper.ED25519
		case 'B':
			keyType = helper.ED448
		default:
			return false, errors.New("unsupported key type")
		}
	} else {
		return false, errors.New("public key is empty")
	}

	switch keyType {
	case helper.ED25519:
		if len(publicKeyByte) != ed25519.PublicKeySize {
			return false, errors.New("invalid public key length for ED25519")
		}
		valid := ed25519.Verify(publicKeyByte, payload, signature)
		if !valid {
			return false, errors.New("signature verification failed")
		}
		return true, nil

	case helper.ED448:
		if len(publicKeyByte) != ed448.PublicKeySize {
			return false, errors.New("invalid public key length for ED448")
		}
		verified := ed448.Verify(publicKeyByte, payload, signature, "")
		if !verified {
			return false, errors.New("ED448: signature verification failed")
		}
		return true, nil

	default:
		return false, errors.New("unsupported key type")
	}
}

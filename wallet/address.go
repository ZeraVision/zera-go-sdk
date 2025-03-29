package wallet

import (
	"errors"

	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/zeebo/blake3"
	"golang.org/x/crypto/sha3"
)

// HashPublicKey hashes a public key using the specified algorithm.
func GetWalletAddress(publicKey []byte, hashAlg helper.HashType, keyType helper.KeyType) ([]byte, string, error) {
	var byteAddr []byte
	var processedPublicKey []byte

	switch hashAlg {
	case helper.BLAKE3:
		hasher := blake3.New()
		hasher.Write(publicKey)
		byteAddr = hasher.Sum(nil)
		processedPublicKey = append([]byte("c_"), publicKey...)
	case helper.SHA3_256:
		hasher := sha3.New256()
		hasher.Write(publicKey)
		byteAddr = hasher.Sum(nil)
		processedPublicKey = append([]byte("a_"), publicKey...)
	case helper.SHA3_512:
		hasher := sha3.New512()
		hasher.Write(publicKey)
		byteAddr = hasher.Sum(nil)
		processedPublicKey = append([]byte("b_"), publicKey...)
	default:
		return nil, "", errors.New("unsupported hash algorithm")
	}

	if keyType == helper.ED25519 {
		processedPublicKey = append([]byte("A_"), processedPublicKey...)
	} else if keyType == helper.ED448 {
		processedPublicKey = append([]byte("B_"), processedPublicKey...)
	} else {
		return nil, "", errors.New("unsupported key type")
	}

	return processedPublicKey, transcode.Base58Encode(byteAddr), nil
}

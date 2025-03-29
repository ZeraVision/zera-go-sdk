package transcode

import (
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"strings"
)

func Base58Decode(encoded string) ([]byte, error) {
	decoded := big.NewInt(0)
	base := big.NewInt(58)
	alphabet := make(map[byte]int64)
	for i, char := range base58Alphabet {
		alphabet[byte(char)] = int64(i)
	}

	for i := 0; i < len(encoded); i++ {
		value, ok := alphabet[encoded[i]]
		if !ok {
			return nil, fmt.Errorf("invalid character in input")
		}

		decoded.Mul(decoded, base)
		decoded.Add(decoded, big.NewInt(value))
	}

	decodedBytes := decoded.Bytes()

	for i := 0; i < len(encoded); i++ {
		if encoded[i] == base58Alphabet[0] {
			decodedBytes = append([]byte{0x00}, decodedBytes...)
		} else {
			break
		}
	}

	return decodedBytes, nil
}

func HexDecode(encoded string) ([]byte, error) {
	decoded, err := hex.DecodeString(encoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex: %w", err)
	}
	return decoded, nil
}

func HashToHexByte(stringHash string) ([]byte, error) {
	var transactionHashByte []byte
	var err error

	re := regexp.MustCompile(`(.*)(s\d+)$`)
	matches := re.FindStringSubmatch(stringHash)

	if len(matches) > 2 {
		result, err := HexDecode(matches[1])
		if err != nil {
			log.Printf("Error decoding hex: %v", err)
		}

		transactionHashByte = append([]byte(result), []byte(matches[2])...)

	} else {

		transactionHashByte, err = HexDecode(stringHash)
		if err != nil {
			log.Printf("Error decoding hex: %v", err)
		}

	}

	return transactionHashByte, err
}

// * Note this only works for single keys
func Base58DecodePublicKey(publicKey string) (prefix []byte, public []byte, combined []byte, err error) {
	// Find the last occurrence of '_'
	lastUnderscoreIndex := strings.LastIndex(publicKey, "_")

	// If no underscore is found, just decode everything
	if lastUnderscoreIndex == -1 {
		result, err := Base58Decode(publicKey)

		if err != nil {
			fmt.Println("Base58DecodePublicKey: " + err.Error())
		}

		return nil, result, result, err
	}

	// Extract the prefix and the part to decode
	prefixStr := publicKey[:lastUnderscoreIndex+1]
	toDecode := publicKey[lastUnderscoreIndex+1:]

	// Decode the part after the underscore
	decoded, err := Base58Decode(toDecode)
	if err != nil {
		return nil, nil, nil, err
	}

	// Prepend the prefix to the decoded byte array
	decodedWithPrefix := append([]byte(prefixStr), decoded...)

	return []byte(prefixStr), decoded, decodedWithPrefix, err
}

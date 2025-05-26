package helper

import (
	"fmt"
	"strings"
)

func DetermineKeyType(publicKeyBase58 string) (KeyType, error) {

	publicKeyBase58 = strings.TrimPrefix(publicKeyBase58, "r_")

	keyLetter := strings.Split(publicKeyBase58, "_")[0]
	switch keyLetter {
	case "A":
		return ED25519, nil
	case "B":
		return ED448, nil
	default:
		return 0, fmt.Errorf("unknown key type for public key: %s", publicKeyBase58)
	}
}

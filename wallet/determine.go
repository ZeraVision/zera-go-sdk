package wallet

import (
	"log"
	"strings"

	"github.com/ZeraVision/zera-go-sdk/helper"
)

func DetermineKeyType(publicKeyBase58 string) (helper.KeyType, error) {

	if strings.HasPrefix(publicKeyBase58, "r_") {
		publicKeyBase58 = publicKeyBase58[2:]
	}

	var keyType helper.KeyType
	if strings.HasPrefix(publicKeyBase58, "A") {
		keyType = helper.ED25519
	} else if strings.HasPrefix(publicKeyBase58, "B") {
		keyType = helper.ED448
	} else if strings.HasPrefix(publicKeyBase58, "special") {
		log.Printf("Special key type inferred: %s", publicKeyBase58)
		keyType = helper.SPECIAL
	}

	return keyType, nil
}

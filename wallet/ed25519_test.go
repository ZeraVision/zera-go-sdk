package wallet_test

import (
	"testing"

	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/wallet"
)

func TestGenerateEd25519_BLAKE3(t *testing.T) {
	testGenerateEd25519(t, helper.BLAKE3, "A_c_7rVZgQm5uPUvkpJDXo9L7hrPA2gcev9G2aKfDNZiKAvx", "23ULwo87vyjEUtZZyKjgVa34b3VE7d6kfWY9MBDM8nVb")
}

func TestGenerateEd25519_SHA3_256(t *testing.T) {
	testGenerateEd25519(t, helper.SHA3_256, "A_a_7rVZgQm5uPUvkpJDXo9L7hrPA2gcev9G2aKfDNZiKAvx", "6GmjroVZfYGzZFGqBh4vp4eypYDMB6Doi4uczvyUMJaT")
}

func TestGenerateEd25519_SHA3_512(t *testing.T) {
	testGenerateEd25519(t, helper.SHA3_512, "A_b_7rVZgQm5uPUvkpJDXo9L7hrPA2gcev9G2aKfDNZiKAvx", "36yjBix8SWW89NFQTaaQJ62dWAoymaJgViKnyityMtSRHm6UMw7uU997WJCXMgoTZQGDQDpCGF44PXyEtgyiQRGR")
}

func testGenerateEd25519(t *testing.T, hashAlgorithm helper.HashType, expectedPublicKey, expectedAddress string) {
	mnemonic := "crumble tattoo grape hurry pizza inject remind play believe museum thing mosquito"
	expectedPrivateKey := "2fweChsECpmRDP5yFgJjoJi2PugxmGAoZcmYkr9kDqNKJGAUzm3DwCPKUq8ZTct3occco1reG2fishuDNchQF9vU"

	//mnemonic = "" //! specifying empty mnemonic will generate random entropy non-BIP39 based

	keyType := helper.ED25519

	b58PrivateKey, b58PublicKey, b58Address, err := wallet.GenerateEd25519(mnemonic, hashAlgorithm, keyType)
	if err != nil {
		t.Fatalf("Error generating key pair: %v", err)
	}

	if b58PrivateKey != expectedPrivateKey {
		t.Errorf("Private Key (Base58) mismatch. Expected: %s, Got: %s", expectedPrivateKey, b58PrivateKey)
	}

	if b58PublicKey != expectedPublicKey {
		t.Errorf("Public Key (Base58) mismatch. Expected: %s, Got: %s", expectedPublicKey, b58PublicKey)
	}

	if b58Address != expectedAddress {
		t.Errorf("Address (B58) mismatch. Expected: %s, Got: %s", expectedAddress, b58Address)
	}
}

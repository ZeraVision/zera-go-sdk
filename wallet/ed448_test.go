package wallet_test

import (
	"testing"

	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/wallet"
)

func TestGenerateEd448_BLAKE3(t *testing.T) {
	testGenerateEd448(t, helper.BLAKE3, "B_c_7epXSHxgXp6e3ogtddX68dZ7Ez6sa5xoKXV9UntbK65CB7tfCNEPV6U61bBaEcviCSFsXH6Cdr3rRD", "AJKR2m2yepRZYLKKABGRuzM2ihnYmxcftkBiLdBoc9ix")
}

func TestGenerateEd448_SHA3_256(t *testing.T) {
	testGenerateEd448(t, helper.SHA3_256, "B_a_7epXSHxgXp6e3ogtddX68dZ7Ez6sa5xoKXV9UntbK65CB7tfCNEPV6U61bBaEcviCSFsXH6Cdr3rRD", "BUk5fdV288X7R9LFc5mgDkiEjt2kb5EMdQRESVrD4SCD")
}

func TestGenerateEd448_SHA3_512(t *testing.T) {
	testGenerateEd448(t, helper.SHA3_512, "B_b_7epXSHxgXp6e3ogtddX68dZ7Ez6sa5xoKXV9UntbK65CB7tfCNEPV6U61bBaEcviCSFsXH6Cdr3rRD", "2AUPfbGu5YcUFgRbSFmovnNvLkTxwe36xkbLAVa7rBjiq2bLk1sjjDCraauaJmV19V8gd4jMBXNBuvHrLMzskyBx")
}

func testGenerateEd448(t *testing.T, hashAlgorithm helper.HashType, expectedPublicKey, expectedAddress string) {
	mnemonic := "crumble tattoo grape hurry pizza inject remind play believe museum thing mosquito"
	expectedPrivateKey := "DzRXgQkou2SQKcVY8enGwmhRYeudKeCJV6gKQ73o5BWpg7sGh4SmhidPtax9KigAuGYUctAYfSKS9L"

	keyType := helper.ED448

	b58PrivateKey, b58PublicKey, b58Address, err := wallet.GenerateEd448(mnemonic, hashAlgorithm, keyType)
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

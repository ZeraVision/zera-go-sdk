package wallet_test

import (
	"testing"

	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/wallet"
)

func TestGenerateEd25519_BLAKE3(t *testing.T) {
	testGenerateEd25519(t, helper.BLAKE3, "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7", "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR")
}

func TestGenerateEd25519_SHA3_256(t *testing.T) {
	testGenerateEd25519(t, helper.SHA3_256, "A_a_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7", "QK2KwEe1qKng1mzfiyDaQMKqYzFvman5CPdEVyRy1PV")
}

func TestGenerateEd25519_SHA3_512(t *testing.T) {
	testGenerateEd25519(t, helper.SHA3_512, "A_b_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7", "2hzpMgngf5zW6QMuQePVdtrMqdYNMC6mdBaWS7S458rRFUPTSwSXgwKMGVfEDuNejR5nWTua7evAyNi48ptNgbmR")
}

func testGenerateEd25519(t *testing.T, hashAlgorithm helper.HashType, expectedPublicKey, expectedAddress string) {
	mnemonic := "crumble tattoo grape hurry pizza inject remind play believe museum thing mosquito"
	expectedPrivateKey := "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs"

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

package wallet_test

import (
	"testing"

	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/wallet"
)

func TestGenerateEd448_BLAKE3(t *testing.T) {
	testGenerateEd448(t, helper.BLAKE3, "B_c_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV", "Hv3KUwrmR8C8XVSxuJFJrQqeDixeDnakUTkUUMZkFCUS")
}

func TestGenerateEd448_SHA3_256(t *testing.T) {
	testGenerateEd448(t, helper.SHA3_256, "B_a_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV", "GmAMncQSf9xxCcyib1Xx7jVdXesD868s86XJiwTTspU1")
}

func TestGenerateEd448_SHA3_512(t *testing.T) {
	testGenerateEd448(t, helper.SHA3_512, "B_b_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV", "4RsD3dH143UCfxscj8xPLGHJCaWeWZmEVUoF42VbkNss5v2UmojFFqbtZ9M5qwPH94euxq3bHodfM2qSiYnnodtU")
}

func testGenerateEd448(t *testing.T, hashAlgorithm helper.HashType, expectedPublicKey, expectedAddress string) {
	mnemonic := "crumble tattoo grape hurry pizza inject remind play believe museum thing mosquito"
	expectedPrivateKey := "HYkGjJY8hjEAxLe1UFzEni5mANwbvTquvTV6mgMT6Qp2Ee1CFYC8tVNfdqyJ9ZwnwsYRUwfMg15suW"

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

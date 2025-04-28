package helper_test

import (
	"math/big"
	"testing"

	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/parts"
	"github.com/ZeraVision/zera-go-sdk/transfer"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/proto"
)

func init() {
	godotenv.Load("../.env")
}

func TestEd25519(t *testing.T) {
	testSignature(t, "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR", "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7", "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs", helper.ED25519)
}

func TestEd448(t *testing.T) {
	testSignature(t, "Hv3KUwrmR8C8XVSxuJFJrQqeDixeDnakUTkUUMZkFCUS", "B_c_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV", "HYkGjJY8hjEAxLe1UFzEni5mANwbvTquvTV6mgMT6Qp2Ee1CFYC8tVNfdqyJ9ZwnwsYRUwfMg15suW", helper.ED448)
}

func testSignature(t *testing.T, address, testPublic, testPrivate string, keyType helper.KeyType) {

	inputs := []transfer.Inputs{}

	inputs = append(inputs, transfer.Inputs{
		B58Address: address,
		KeyType:    keyType,
		PublicKey:  testPublic,
		PrivateKey: testPrivate,
		Amount:     "1.01",
		FeePercent: 100,
	})

	outputs := map[string]string{}

	outputs["outputAddr1"] = "1.01"

	baseFeeID := "$ZRA+0000"
	baseFeeAmountParts := "1000000000" // 1 zra

	nonceInfo := nonce.NonceInfo{
		Override: []uint64{5}, // to test the signature this doesnt need to be valid. See nonce_test for usage
	}

	partsInfo := parts.PartsInfo{
		Symbol:   "$ZRA+0000",            //* symbol specified here...
		Override: big.NewInt(1000000000), // to test the signature this doesnt need to be valid. See nonce_test for usage
	}

	// via indexer
	txn, err := transfer.CreateCoinTxn(nonceInfo, partsInfo, inputs, outputs, baseFeeID, baseFeeAmountParts, nil, nil, 5)
	// via validator
	//txn, err := transfer.CreateCoinTxn(false, inputs, outputs, testvars.TEST_GRPC_ADDR, "", symbol, baseFeeID, baseFeeAmountParts, nil, nil)

	if err != nil {
		t.Errorf("Error creating transaction: %s", err)
	}

	// Grab signature
	signature := txn.Auth.Signature[0]

	// Remove signature & hash before verification
	txn.Auth.Signature = nil
	txn.Base.Hash = nil

	txnBytes, err := proto.Marshal(txn)
	if err != nil {
		t.Errorf("Error marshalling transaction: %s", err)
	}

	ok, err := helper.Verify(testPublic, txnBytes, signature)

	if err != nil {
		t.Errorf("Error verifying signature: %s", err)
	}

	if !ok {
		t.Errorf("Signature verification failed")
	}
}

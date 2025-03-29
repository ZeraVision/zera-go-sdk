package transfer_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/transfer"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func Test25519OnetoOne(t *testing.T) {
	inputs := []transfer.Inputs{
		{
			B58Address:         "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR",
			KeyType:            helper.ED25519,
			PublicKey:          "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7",
			PrivateKey:         "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs",
			Amount:             big.NewFloat(1.23456),
			FeePercent:         100,
			ContractFeePercent: nil,
		},
	}
	outputs := map[string]*big.Float{
		"b58addr1": big.NewFloat(1.23456),
	}

	testCoin(t, inputs, outputs, "$ZRA+0000", "$ZRA+0000", "1000000000")
}

func Test25519OnetoMany(t *testing.T) {
	inputs := []transfer.Inputs{
		{
			B58Address:         "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR",
			KeyType:            helper.ED25519,
			PublicKey:          "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7",
			PrivateKey:         "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs",
			Amount:             big.NewFloat(1.23456),
			FeePercent:         100,
			ContractFeePercent: nil,
		},
	}
	outputs := map[string]*big.Float{
		"b58addr1": big.NewFloat(1),
		"b58addr2": big.NewFloat(0.23456),
	}

	testCoin(t, inputs, outputs, "$ZRA+0000", "$ZRA+0000", "1000000000")
}

func Test448OnetoOne(t *testing.T) {
	inputs := []transfer.Inputs{
		{
			B58Address:         "Hv3KUwrmR8C8XVSxuJFJrQqeDixeDnakUTkUUMZkFCUS",
			KeyType:            helper.ED448,
			PublicKey:          "B_c_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV",
			PrivateKey:         "HYkGjJY8hjEAxLe1UFzEni5mANwbvTquvTV6mgMT6Qp2Ee1CFYC8tVNfdqyJ9ZwnwsYRUwfMg15suW",
			Amount:             big.NewFloat(1.23456),
			FeePercent:         100,
			ContractFeePercent: nil,
		},
	}
	outputs := map[string]*big.Float{
		"b58addr1": big.NewFloat(1.23456),
	}

	testCoin(t, inputs, outputs, "$ZRA+0000", "$ZRA+0000", "1000000000")
}

func Test448OnetoMany(t *testing.T) {
	inputs := []transfer.Inputs{
		{
			B58Address:         "Hv3KUwrmR8C8XVSxuJFJrQqeDixeDnakUTkUUMZkFCUS",
			KeyType:            helper.ED448,
			PublicKey:          "B_c_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV",
			PrivateKey:         "HYkGjJY8hjEAxLe1UFzEni5mANwbvTquvTV6mgMT6Qp2Ee1CFYC8tVNfdqyJ9ZwnwsYRUwfMg15suW",
			Amount:             big.NewFloat(1.23456),
			FeePercent:         100,
			ContractFeePercent: nil,
		},
	}
	outputs := map[string]*big.Float{
		"b58addr1": big.NewFloat(1),
		"b58addr2": big.NewFloat(0.23456),
	}

	testCoin(t, inputs, outputs, "$ZRA+0000", "$ZRA+0000", "1000000000")
}

func TestManytoOne(t *testing.T) {
	inputs := []transfer.Inputs{
		{
			B58Address:         "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR",
			KeyType:            helper.ED25519,
			PublicKey:          "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7",
			PrivateKey:         "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs",
			Amount:             big.NewFloat(1.23456),
			FeePercent:         50,
			ContractFeePercent: nil,
		},
		{
			B58Address:         "Hv3KUwrmR8C8XVSxuJFJrQqeDixeDnakUTkUUMZkFCUS",
			KeyType:            helper.ED448,
			PublicKey:          "B_c_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV",
			PrivateKey:         "HYkGjJY8hjEAxLe1UFzEni5mANwbvTquvTV6mgMT6Qp2Ee1CFYC8tVNfdqyJ9ZwnwsYRUwfMg15suW",
			Amount:             big.NewFloat(1.23456),
			FeePercent:         50,
			ContractFeePercent: nil,
		},
	}
	outputs := map[string]*big.Float{
		"b58addr1": big.NewFloat(2.46912),
	}

	testCoin(t, inputs, outputs, "$ZRA+0000", "$ZRA+0000", "1000000000")
}

func TestManytoMany(t *testing.T) {
	inputs := []transfer.Inputs{
		{
			B58Address:         "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR",
			KeyType:            helper.ED25519,
			PublicKey:          "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7",
			PrivateKey:         "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs",
			Amount:             big.NewFloat(1.23456),
			FeePercent:         50,
			ContractFeePercent: nil,
		},
		{
			B58Address:         "Hv3KUwrmR8C8XVSxuJFJrQqeDixeDnakUTkUUMZkFCUS",
			KeyType:            helper.ED448,
			PublicKey:          "B_c_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV",
			PrivateKey:         "HYkGjJY8hjEAxLe1UFzEni5mANwbvTquvTV6mgMT6Qp2Ee1CFYC8tVNfdqyJ9ZwnwsYRUwfMg15suW",
			Amount:             big.NewFloat(1.23456),
			FeePercent:         50,
			ContractFeePercent: nil,
		},
	}
	outputs := map[string]*big.Float{
		"b58addr1": big.NewFloat(2.00),
		"b58addr2": big.NewFloat(0.46912),
	}

	testCoin(t, inputs, outputs, "$ZRA+0000", "$ZRA+0000", "1000000000")
}

func testCoin(t *testing.T, inputs []transfer.Inputs, outputs map[string]*big.Float, symbol, baseFeeID, baseFeeAmountParts string) {
	// Indexer
	txn, err := transfer.CreateCoinTxn(true, inputs, outputs, "https://indexer.zera.vision", os.Getenv("INDEXER_API_KEY"), symbol, baseFeeID, baseFeeAmountParts, nil, nil)

	// Validator
	grpcAddr := "routing.zera.vision"
	//txn, err := transfer.CreateCoinTxn(false, inputs, outputs, grpcAddr+ ":50051", "", symbol, baseFeeID, baseFeeAmountParts, nil, nil)

	if err != nil {
		t.Errorf("Error creating transaction: %s", err)
	}

	_, err = transfer.SendCoinTXN(grpcAddr+":50052", txn)

	if err != nil {
		t.Errorf("Error sending transaction: %s", err)
	}
}

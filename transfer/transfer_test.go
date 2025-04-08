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
			B58Address:         "23ULwo87vyjEUtZZyKjgVa34b3VE7d6kfWY9MBDM8nVb",
			KeyType:            helper.ED25519,
			PublicKey:          "A_c_7rVZgQm5uPUvkpJDXo9L7hrPA2gcev9G2aKfDNZiKAvx",
			PrivateKey:         "2fweChsECpmRDP5yFgJjoJi2PugxmGAoZcmYkr9kDqNKJGAUzm3DwCPKUq8ZTct3occco1reG2fishuDNchQF9vU",
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
			B58Address:         "23ULwo87vyjEUtZZyKjgVa34b3VE7d6kfWY9MBDM8nVb",
			KeyType:            helper.ED25519,
			PublicKey:          "A_c_7rVZgQm5uPUvkpJDXo9L7hrPA2gcev9G2aKfDNZiKAvx",
			PrivateKey:         "2fweChsECpmRDP5yFgJjoJi2PugxmGAoZcmYkr9kDqNKJGAUzm3DwCPKUq8ZTct3occco1reG2fishuDNchQF9vU",
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
			B58Address:         "AJKR2m2yepRZYLKKABGRuzM2ihnYmxcftkBiLdBoc9ix",
			KeyType:            helper.ED448,
			PublicKey:          "B_c_7epXSHxgXp6e3ogtddX68dZ7Ez6sa5xoKXV9UntbK65CB7tfCNEPV6U61bBaEcviCSFsXH6Cdr3rRD",
			PrivateKey:         "DzRXgQkou2SQKcVY8enGwmhRYeudKeCJV6gKQ73o5BWpg7sGh4SmhidPtax9KigAuGYUctAYfSKS9L",
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
			B58Address:         "AJKR2m2yepRZYLKKABGRuzM2ihnYmxcftkBiLdBoc9ix",
			KeyType:            helper.ED448,
			PublicKey:          "B_c_7epXSHxgXp6e3ogtddX68dZ7Ez6sa5xoKXV9UntbK65CB7tfCNEPV6U61bBaEcviCSFsXH6Cdr3rRD",
			PrivateKey:         "DzRXgQkou2SQKcVY8enGwmhRYeudKeCJV6gKQ73o5BWpg7sGh4SmhidPtax9KigAuGYUctAYfSKS9L",
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
			B58Address:         "23ULwo87vyjEUtZZyKjgVa34b3VE7d6kfWY9MBDM8nVb",
			KeyType:            helper.ED25519,
			PublicKey:          "A_c_7rVZgQm5uPUvkpJDXo9L7hrPA2gcev9G2aKfDNZiKAvx",
			PrivateKey:         "2fweChsECpmRDP5yFgJjoJi2PugxmGAoZcmYkr9kDqNKJGAUzm3DwCPKUq8ZTct3occco1reG2fishuDNchQF9vU",
			Amount:             big.NewFloat(1.23456),
			FeePercent:         50,
			ContractFeePercent: nil,
		},
		{
			B58Address:         "AJKR2m2yepRZYLKKABGRuzM2ihnYmxcftkBiLdBoc9ix",
			KeyType:            helper.ED448,
			PublicKey:          "B_c_7epXSHxgXp6e3ogtddX68dZ7Ez6sa5xoKXV9UntbK65CB7tfCNEPV6U61bBaEcviCSFsXH6Cdr3rRD",
			PrivateKey:         "DzRXgQkou2SQKcVY8enGwmhRYeudKeCJV6gKQ73o5BWpg7sGh4SmhidPtax9KigAuGYUctAYfSKS9L",
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
			B58Address:         "23ULwo87vyjEUtZZyKjgVa34b3VE7d6kfWY9MBDM8nVb",
			KeyType:            helper.ED25519,
			PublicKey:          "A_c_7rVZgQm5uPUvkpJDXo9L7hrPA2gcev9G2aKfDNZiKAvx",
			PrivateKey:         "2fweChsECpmRDP5yFgJjoJi2PugxmGAoZcmYkr9kDqNKJGAUzm3DwCPKUq8ZTct3occco1reG2fishuDNchQF9vU",
			Amount:             big.NewFloat(1.23456),
			FeePercent:         50,
			ContractFeePercent: nil,
		},
		{
			B58Address:         "AJKR2m2yepRZYLKKABGRuzM2ihnYmxcftkBiLdBoc9ix",
			KeyType:            helper.ED448,
			PublicKey:          "B_c_7epXSHxgXp6e3ogtddX68dZ7Ez6sa5xoKXV9UntbK65CB7tfCNEPV6U61bBaEcviCSFsXH6Cdr3rRD",
			PrivateKey:         "DzRXgQkou2SQKcVY8enGwmhRYeudKeCJV6gKQ73o5BWpg7sGh4SmhidPtax9KigAuGYUctAYfSKS9L",
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

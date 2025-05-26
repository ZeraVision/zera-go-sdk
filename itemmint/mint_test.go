package itemmint_test

import (
	"math/big"
	"testing"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/itemmint"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/testvars"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestItemMint(t *testing.T) {
	mintFromAddr := "E9U2FT5gRkqKS3MsDWYVpVmvLX7XaC7JGL38pZSpYEn3"
	publicKey := "r_A_a_6E4RM19bfUJHBKcoZCvfiEmZ7JVWpgbFGgtadwjH9Xoo"
	privateKey := "31WUV7cUHUVWtoCUtf6WqTHooMCQ4Cgsikt6j98QFbBaD3cYShXWeNXFSScCJwYRxNj8ismxrBU6LiFPL51NBZXf"

	// Nonce
	nonceReq, err := nonce.MakeNonceRequest(mintFromAddr)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	nonceInfo := nonce.NonceInfo{
		UseIndexer:    false,
		NonceReqs:     []*pb.NonceRequest{nonceReq},
		ValidatorAddr: testvars.TEST_GRPC_ADDRESS,
	}

	// NFT mint parameters
	itemToAddr := "CY7JXLDTwfqYUZsLJ58bHScFpK2gSgFLQy2CfAVGDBBt" // recipient
	contractId := "$NFT+0000"
	itemId := big.NewInt(1) // should be unique per NFT
	baseFeeSymbol := "$ZRA+0000"
	baseFeeParts := "10000000000"

	// Example metadata parameters
	parameters := []*pb.KeyValuePair{
		{Key: "name", Value: "Test NFT"},
		{Key: "description", Value: "A test NFT for unit testing."},
	}

	// Example voting weight
	votingWeight := new(big.Int)
	votingWeight = big.NewInt(1)

	// Example contract fees (use nil or a real value if available)
	var contractFees *pb.ItemContractFees = nil // Build item contract fees function can assist with this

	// Example expiry and validFrom (use nil or a real value if needed)
	var expiry *uint64 = nil
	var validFrom *uint64 = nil

	txn, err := itemmint.CreateItemMintTxn(
		nonceInfo,
		contractId,
		itemId,
		itemToAddr,
		publicKey,
		privateKey,
		baseFeeSymbol,
		baseFeeParts,
		parameters,
		expiry,
		validFrom,
		votingWeight,
		contractFees,
	)

	if err != nil {
		t.Errorf("Error creating item mint transaction: %s", err)
	}

	_, err = itemmint.SendItemMintTXN(testvars.TEST_GRPC_ADDRESS+":50052", txn)

	if err != nil {
		t.Errorf("Error sending item mint transaction: %s", err)
	}

	if testvars.TEST_GRPC_ADDRESS == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

func TestIndexerMint(t *testing.T) {
	mintFromAddr := "BjDzM8uhJQgzJ2xNQuGmQLHk7TtTDbYM9uGtae8yerka"
	publicKey := "r_A_c_ApWitx3UxKmMKr8H1QdqrKU347zgJ5SN62Xax9fyWaki"
	privateKey := "39mNSLF3XBs8dSvcD5XsM6F1zkopcB76rybpycuFPRKNMEUH4k4QM2NChc9xRAGUMUgtJfVgLZud99GhTFfszx4a"

	// Nonce
	nonceReq, err := nonce.MakeNonceRequest(mintFromAddr)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	nonceInfo := nonce.NonceInfo{
		UseIndexer:    false,
		NonceReqs:     []*pb.NonceRequest{nonceReq},
		ValidatorAddr: testvars.TEST_GRPC_ADDRESS,
	}

	// NFT mint parameters
	itemToAddr := "5dFWUxkYqYEE6cRcoxa9fXR9v9VLtqJ9ebpeMyZSPD72" // ZV Main
	contractId := "$INDEXERIP+0000"
	itemId := big.NewInt(0) // should be unique per NFT
	baseFeeSymbol := "$ZRA+0000"
	baseFeeParts := "10000000000"

	// Example contract fees (use nil or a real value if available)
	var contractFees *pb.ItemContractFees = nil // Build item contract fees function can assist with this

	// Example expiry and validFrom (use nil or a real value if needed)
	var expiry *uint64 = nil
	var validFrom *uint64 = nil

	txn, err := itemmint.CreateItemMintTxn(
		nonceInfo,
		contractId,
		itemId,
		itemToAddr,
		publicKey,
		privateKey,
		baseFeeSymbol,
		baseFeeParts,
		nil,
		expiry,
		validFrom,
		nil,
		contractFees,
	)

	if err != nil {
		t.Errorf("Error creating item mint transaction: %s", err)
	}

	_, err = itemmint.SendItemMintTXN(testvars.TEST_GRPC_ADDRESS+":50052", txn)

	if err != nil {
		t.Errorf("Error sending item mint transaction: %s", err)
	}

	if testvars.TEST_GRPC_ADDRESS == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

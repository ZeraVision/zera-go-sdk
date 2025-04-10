package expenseratio_test

import (
	"testing"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/expenseratio"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/testvars"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestExpenseRatio(t *testing.T) {
	fromAddr := "E9U2FT5gRkqKS3MsDWYVpVmvLX7XaC7JGL38pZSpYEn3" // in this test case, sending results to self
	publicKey := "r_A_a_6E4RM19bfUJHBKcoZCvfiEmZ7JVWpgbFGgtadwjH9Xoo"
	privateKey := "31WUV7cUHUVWtoCUtf6WqTHooMCQ4Cgsikt6j98QFbBaD3cYShXWeNXFSScCJwYRxNj8ismxrBU6LiFPL51NBZXf"

	// Indexer
	// nonceInfo := nonce.NonceInfo{
	// 	UseIndexer:    true,
	// 	Addresses:       []string{mintFromAddr},
	// 	IndexerURL:    "https://indexer.zera.vision",
	// 	Authorization: os.Getenv("INDEXER_API_KEY"),
	// }

	nonceReq, err := nonce.MakeNonceRequest(fromAddr)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validator
	nonceInfo := nonce.NonceInfo{
		UseIndexer:    false,
		NonceReqs:     []*pb.NonceRequest{nonceReq},
		ValidatorAddr: testvars.TEST_GRPC_ADDRESS,
	}

	symbol := "$BENCHY+0000"
	calledAddrs := []string{
		"addr1",
		"addr2",
	}

	baseFeeSymbol := symbol      // some ACE authorized token
	baseFeeParts := "1000000000" // a number of parts sufficient to cover the fee with your symbol

	txn, err := expenseratio.ExpenseRatioTxn(nonceInfo, symbol, calledAddrs, fromAddr, publicKey, privateKey, baseFeeSymbol, baseFeeParts) // feeID and feeAmountParts are not used in this test case

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = expenseratio.SendExpenseRatioTXN(testvars.TEST_GRPC_ADDRESS+":50052", txn)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if testvars.TEST_GRPC_ADDRESS == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

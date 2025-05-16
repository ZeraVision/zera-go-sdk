package allowance_test

import (
	"math/big"
	"os"
	"testing"
	"time"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/allowance"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/testvars"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestAllowance(t *testing.T) {
	fromAddr := "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR"
	publicKey := "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7"
	privateKey := "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs"

	fromAddr = os.Getenv("TEST_ADDR")
	publicKey = os.Getenv("TEST_PUBLIC")
	privateKey = os.Getenv("TEST_PRIVATE")

	// Indexer
	// nonceInfo := nonce.NonceInfo{
	// 	UseIndexer:    true,
	// 	Addresses:       []string{fromAddr},
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

	symbol := "$ZRA+0000"

	baseFeeSymbol := "$ZRA+0000" // some ACE authorized token
	baseFeeParts := "1000000000" // a number of parts sufficient to cover the fee with your symbol

	var currencyEquivalent *float64
	var amount *big.Int
	if true {
		currencyEquivalentV := 1.234 // $1.234
		currencyEquivalent = &currencyEquivalentV
	} else {
		amount = big.NewInt(100_000_000_000) // 100 billion parts
	}

	var periodMonths *uint32
	var periodSeconds *uint32
	if true {
		periodMonthsV := uint32(1) // 1 month(s)
		periodMonths = &periodMonthsV
	} else {
		periodSecondsV := uint32(60) // 60 seconds
		periodSeconds = &periodSecondsV
	}

	authorize := false

	allowanceDetails := allowance.AllowanceDetails{
		Authorize:          authorize,
		WalletAddr:         "QK2KwEe1qKng1mzfiyDaQMKqYzFvman5CPdEVyRy1PV",
		CurrencyEquivalent: currencyEquivalent,
		Amount:             amount,
		PeriodMonths:       periodMonths,
		PeriodSeconds:      periodSeconds,
		StartTime:          time.Now().Unix(), // update as needed
	}

	txn, err := allowance.CreateAllowanceTxn(nonceInfo, symbol, allowanceDetails, publicKey, privateKey, baseFeeSymbol, baseFeeParts)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = allowance.SendAllowanceTxn(testvars.TEST_GRPC_ADDRESS+":50052", txn)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if testvars.TEST_GRPC_ADDRESS == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

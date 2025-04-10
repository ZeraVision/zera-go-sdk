package currencyequivalent_test

import (
	"math/big"
	"os"
	"testing"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	ace "github.com/ZeraVision/zera-go-sdk/currencyequivalent"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/testvars"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestAce(t *testing.T) {

	aceAddr := os.Getenv("ACE_ADDR")
	privateKey := os.Getenv("ACE_PRIVATE_KEY")
	publicKey := os.Getenv("ACE_PUBLIC_KEY")

	// Indexer
	// nonceInfo := nonce.NonceInfo{
	// 	UseIndexer:    true,
	// 	Address:       aceAddr,
	// 	IndexerURL:    "https://indexer.zera.vision",
	// 	Authorization: os.Getenv("INDEXER_API_KEY"),
	// }

	nonceReq, err := nonce.MakeNonceRequest(aceAddr)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validator
	nonceInfo := nonce.NonceInfo{
		UseIndexer:    false,
		NonceReqs:     []*pb.NonceRequest{nonceReq},
		ValidatorAddr: testvars.TEST_GRPC_ADDRESS,
	}

	rate := big.NewFloat(123.456)                                                    // usd
	networkRateString := rate.Mul(rate, big.NewFloat(1e18)).Text('f', 0)             // network takes a scaled 1e18 version of the rate
	authorized := true                                                               // true authorize, false deauthorize
	maxStake := big.NewFloat(123_456.589)                                            // usd
	networkMaxStakeString := maxStake.Mul(maxStake, big.NewFloat(1e18)).Text('f', 0) // network takes a scaled 1e18 version of the max stake

	data := []ace.AceData{
		{
			Symbol:     "$BENCHY+0000",
			Rate:       networkRateString,
			Authorized: &authorized,
			MaxStake:   &networkMaxStakeString,
		},
	}

	txn, err := ace.CreateAceTxn(nonceInfo, data, publicKey, privateKey, "$ZRA+0000", "10000000000")

	if err != nil {
		t.Fatalf("Error creating transaction: %s", err)
	}

	_, err = ace.SendAceTXN(testvars.TEST_GRPC_ADDRESS+":50052", txn)

	if err != nil {
		t.Fatalf("Error sending transaction: %s", err)
	}

	if testvars.TEST_GRPC_ADDRESS == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

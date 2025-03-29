package currencyequivalent_test

import (
	"math/big"
	"testing"

	self "github.com/ZeraVision/zera-go-sdk/currencyequivalent"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestSelfCurrencyEquivalent(t *testing.T) {
	selfAddr := "5EyfBMt2XNTKQuJqjHjngph6oF2Cb8bYJSD7QxaWdMp2VhuD9D2H5p3ZE9mWyp119SnyiMVHkTztzZQ9aZwTt59h"
	privateKey := "3uhKAgqHD6RZBjd8kKxjDZTz1zQb55NhBjaFpW2HesKHNdnKz9qPMZehbeQ4qecfNFM1aZ3buQzNKsg9hBZHAJea"
	publicKey := "r_A_b_82BWAE3iRDBuaxjDyommMfTLs44HthypLVmWnsA1rvTp"

	// Indexer
	// nonceInfo := nonce.NonceInfo{
	// 	UseIndexer:    true,
	// 	Address:       selfAddr,
	// 	IndexerURL:    "https://indexer.zera.vision",
	// 	Authorization: os.Getenv("INDEXER_API_KEY"),
	// }

	grpcAddr := "routing.zera.vision" // Change grpc addr as required
	grpcAddr = "125.253.87.133"       // override

	nonceReq, err := nonce.MakeNonceRequest(selfAddr)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Validator
	nonceInfo := nonce.NonceInfo{
		UseIndexer:    false,
		NonceReq:      nonceReq,
		ValidatorAddr: grpcAddr + ":50051",
	}

	rate := big.NewFloat(1.01)                                           // usd
	networkRateString := rate.Mul(rate, big.NewFloat(1e18)).Text('f', 0) // network takes a scaled 1e18 version of the rate                                                        // true authorize, false deauthorize

	data := []self.SelfData{
		{
			Symbol: "$BENCHY+0000",
			Rate:   networkRateString,
		},
	}

	txn, err := self.CreateSelfCurrencyEquivalentTxn(nonceInfo, data, publicKey, privateKey, "$ZRA+0000", "10000000000")

	if err != nil {
		t.Fatalf("Error creating transaction: %s", err)
	}

	_, err = self.SendSelfCurrencyEquivalentTXN(grpcAddr+":50052", txn)

	if err != nil {
		t.Fatalf("Error sending transaction: %s", err)
	}

	if grpcAddr == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

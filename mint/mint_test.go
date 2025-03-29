package mint_test

import (
	"testing"

	"github.com/ZeraVision/zera-go-sdk/mint"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestMintEd25519(t *testing.T) {
	mintFromAddr := "E9U2FT5gRkqKS3MsDWYVpVmvLX7XaC7JGL38pZSpYEn3"
	publicKey := "r_A_a_6E4RM19bfUJHBKcoZCvfiEmZ7JVWpgbFGgtadwjH9Xoo"
	privateKey := "31WUV7cUHUVWtoCUtf6WqTHooMCQ4Cgsikt6j98QFbBaD3cYShXWeNXFSScCJwYRxNj8ismxrBU6LiFPL51NBZXf"

	testMint(t, mintFromAddr, publicKey, privateKey)
}

func TestMintEd448(t *testing.T) {
	mintFromAddr := "CY7JXLDTwfqYUZsLJ58bHScFpK2gSgFLQy2CfAVGDBBt"
	publicKey := "r_B_a_S9ukYEEup11N1p5P2ViWRBRT5k8HNsJfALQsh7T4hLkQeqb9GxgLXvQtXhqYTKxJBBGvhc6HetXd9Z"
	privateKey := "6qGrAyMom9Q938uDrDvqyeTDQjXNJt4iQgmsNMQqfEXVvj7AbJBW3Rma8vNUYmNHrJqvCEhLq6D5ch"

	testMint(t, mintFromAddr, publicKey, privateKey)
}

func testMint(t *testing.T, mintFromAddr, publicKey, privateKey string) {

	// Indexer
	// nonceInfo := nonce.NonceInfo{
	// 	UseIndexer:    true,
	// 	Address:       mintFromAddr,
	// 	IndexerURL:    "https://indexer.zera.vision",
	// 	Authorization: os.Getenv("INDEXER_API_KEY"),
	// }

	nonceReq, err := nonce.MakeNonceRequest(mintFromAddr)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	grpcAddr := "routing.zera.vision" // Change grpc addr as required
	grpcAddr = "125.253.87.133"       // override

	// Validator
	nonceInfo := nonce.NonceInfo{
		UseIndexer:    false,
		NonceReq:      nonceReq,
		ValidatorAddr: grpcAddr + ":50051",
	}

	// Random addr for demo purposes
	mintToAddr := "CY7JXLDTwfqYUZsLJ58bHScFpK2gSgFLQy2CfAVGDBBt"
	symbol := "$BENCHY+0000"
	amountParts := "10000000000" // full coins is calculated as amountParts / parts per coin (denomination)
	baseFeeSymbol := symbol      // some ACE authorized token
	baseFeeParts := "1000000000" // a number of parts sufficient to cover the fee with your symbol

	txn, err := mint.CreateMintTxn(nonceInfo, symbol, amountParts, mintToAddr, publicKey, privateKey, baseFeeSymbol, baseFeeParts)

	if err != nil {
		t.Errorf("Error creating transaction: %s", err)
	}

	// Change grpc addr as required
	_, err = mint.SendMintTXN(grpcAddr+":50052", txn)

	if err != nil {
		t.Errorf("Error sending transaction: %s", err)
	}

	if grpcAddr == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

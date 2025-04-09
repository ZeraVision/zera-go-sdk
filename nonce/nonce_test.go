package nonce_test

import (
	"os"
	"testing"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/testvars"
	"github.com/joho/godotenv"
)

const NONCE_TEST_ADDR = "48aPY5LHV6rHXAS5ciZNYPTGYV1fm1k4BQ8Wakh2B1xP"

func init() {
	godotenv.Load("../.env")
}

func TestGetNonce_UseIndexer(t *testing.T) {
	nonceInfo := nonce.NonceInfo{
		UseIndexer:    true,
		Addresses:     []string{NONCE_TEST_ADDR},
		IndexerURL:    "https://indexer.zera.vision",
		Authorization: os.Getenv("INDEXER_API_KEY"),
	}

	nonceValue, err := nonce.GetNonce(nonceInfo)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for _, value := range nonceValue {
		if value < 1 {
			t.Fatalf("Expected nonce to be greater than 0, got %d", value)
		}
	}

	t.Logf("Retrieved nonce from Indexer: %d", nonceValue)
}

func TestGetNonce_ValidatorMode(t *testing.T) {
	nonceReq, err := nonce.MakeNonceRequest(NONCE_TEST_ADDR)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	nonceInfo := nonce.NonceInfo{
		UseIndexer:    false,
		NonceReqs:     []*pb.NonceRequest{nonceReq},
		ValidatorAddr: testvars.TEST_GRPC_ADDRESS,
	}

	nonceValue, err := nonce.GetNonce(nonceInfo)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	for _, value := range nonceValue {
		if value < 1 {
			t.Fatalf("Expected nonce to be greater than 0, got %d", value)
		}
	}

	t.Logf("Retrieved nonce from Validator: %d", nonceValue)
}

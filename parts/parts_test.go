package parts_test

import (
	"math/big"
	"os"
	"testing"

	"github.com/ZeraVision/zera-go-sdk/parts"
	"github.com/joho/godotenv"
)

const SAMPLE_SYMBOL = "$ZRA+0000"

func init() {
	godotenv.Load("../.env")
}

func TestGetParts_UseIndexer(t *testing.T) {
	partsInfo := parts.PartsInfo{
		Symbol:        SAMPLE_SYMBOL,
		UseIndexer:    true,
		IndexerUrl:    "https://indexer.zera.vision",
		Authorization: os.Getenv("INDEXER_API_KEY"),
	}

	parts, err := parts.GetParts(partsInfo) // can use api key or bearer token
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if parts == nil {
		t.Fatalf("Expected parts to be non-nil")
	}

	if parts.Cmp(big.NewInt(1)) < 0 {
		t.Fatalf("Expected parts to be greater than or equal to 1, got %s", parts.String())
	}

	t.Logf("Retrieved parts of %s from Indexer: %s", SAMPLE_SYMBOL, parts.String())
}

// not yet possible
func TestGetParts_ValidatorMode(t *testing.T) {
}

package governance_test

import (
	"testing"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/governance"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/testvars"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestVote(t *testing.T) {
	fromAddr := "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR"
	publicKey := "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7"
	privateKey := "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs"

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

	symbol := "$FIBZ+0000"

	var support *bool
	var supportOption *uint32

	// Support / Against
	proposalID := "d88bf4910fe2e480e71c310c91e894230df8950385b10cefb28c033926ab4c9b"
	sBool := true
	support = &sBool

	// Vote option
	// proposalID := "a5ff7558ce466401819cf200b061efa5de8f92102348a628acf2bd5fef37f4a2"
	// voteOption := uint32(1) // lighthouse in this example
	// supportOption = &voteOption

	baseFeeSymbol := "$ZRA+0000" // some ACE authorized token
	baseFeeParts := "1000000000" // a number of parts sufficient to cover the fee with your symbol

	txn, err := governance.CreateVoteTxn(nonceInfo, symbol, proposalID, publicKey, privateKey, baseFeeSymbol, baseFeeParts, support, supportOption)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = governance.SendVoteTxn(testvars.TEST_GRPC_ADDRESS+":50052", txn)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if testvars.TEST_GRPC_ADDRESS == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

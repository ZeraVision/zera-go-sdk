package governance_test

import (
	"testing"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/governance"
	"github.com/ZeraVision/zera-go-sdk/mint"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/testvars"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/joho/godotenv"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func init() {
	godotenv.Load("../.env")
}

func TestProposal(t *testing.T) {
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

	baseFeeSymbol := "$ZRA+0000"  // some ACE authorized token
	baseFeeParts := "10000000000" // a number of parts sufficient to cover the fee with your symbol

	// Configure proposal
	title := "Fibbin Around from https://github.com/ZeraVision/zera-go-sdk/"
	synopsis := "The SDK is a great tool for developers for interact with the ZERA Network!"
	body := `
# Summary
Markdown is seemingly the standard for proposals -- it allows for some good formatting and is easy to read. Apps like the ZERA Vision Explorer work to display various markdown elements.

## Go-SDK
The Go-SDK is a great tool for developers to interact with the ZERA Network. It provides a simple and easy-to-use interface for developers to build applications on top of the ZERA Network and gives examples in test files.

## Associated Transactions
Gimme gimme gimme!!!! -- A mint transaction associated - provided this key is authorized to due it under the rules of the contract, it can happen upon successful "support" result.

# Conclusion
That's it for now -- one dev to the next, I can't wait to see where your creation leads 🔥🔥🔥🔥
	`

	options := []string{} // if you want a multivote option (ie https://explorer.zera.vision/proposal/a5ff7558ce466401819cf200b061efa5de8f92102348a628acf2bd5fef37f4a2, specify objects in this, else specify empty arr)

	// These are only present for adaptive gov type -- most cases wont use these
	var startTimestamp, endTimestamp *timestamppb.Timestamp
	// startTimestamp = timestamppb.Now()
	// endTimestamp = timestamppb.New(startTimestamp.AsTime().Add(24 * 60 * 60 * 1000)) // 24 hours from now

	// If there are transaction(s) associated with it the result of a proposal (only in Support Against will work) then specify them here
	// Ex modified from mint_test.gov
	var testMintTxn *pb.MintTXN
	var testMintBytes []byte
	{
		// Indexer
		// nonceInfo := nonce.NonceInfo{
		// 	UseIndexer:    true,
		// 	Address:       mintFromAddr,
		// 	IndexerURL:    "https://indexer.zera.vision",
		// 	Authorization: os.Getenv("INDEXER_API_KEY"),
		// }

		mintFromKey := "gov_" + symbol // must be contract proposal is being submitted to for chance at being valid
		nonceReq, err := nonce.MakeNonceRequest(mintFromKey)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Validator
		nonceInfo := nonce.NonceInfo{
			UseIndexer:    false,
			NonceReqs:     []*pb.NonceRequest{nonceReq},
			ValidatorAddr: testvars.TEST_GRPC_ADDRESS,
		}

		// Random addr for demo purposes
		mintToAddr := "CY7JXLDTwfqYUZsLJ58bHScFpK2gSgFLQy2CfAVGDBBt"
		symbol := "$FIBZ+0000"
		amountParts := "10000000000" // full coins is calculated as amountParts / parts per coin (denomination)
		baseFeeSymbol := symbol      // some ACE authorized token
		baseFeeParts := "1000000000" // a number of parts sufficient to cover the fee with your symbol

		testMintTxn, err = mint.CreateMintTxn(nonceInfo, symbol, amountParts, mintToAddr, mintFromKey, "", baseFeeSymbol, baseFeeParts)

		if err != nil {
			t.Fatalf("Error creating transaction: %s", err)
		}

		testMintBytes, err = proto.Marshal(testMintTxn)
		if err != nil {
			t.Fatalf("Error serializing transaction: %s", err)
		}
	}

	associatedTxns := []*pb.GovernanceTXN{
		{
			TxnType:       pb.TRANSACTION_TYPE_MINT_TYPE,
			SerializedTxn: testMintBytes,
			TxnHash:       testMintTxn.Base.Hash,
		},
	}

	txn, err := governance.CreateProposalTxn(nonceInfo, symbol, publicKey, privateKey, baseFeeSymbol, baseFeeParts, title, synopsis, body, options, startTimestamp, endTimestamp, associatedTxns)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = governance.SendProposal(testvars.TEST_GRPC_ADDRESS+":50052", txn)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if testvars.TEST_GRPC_ADDRESS == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

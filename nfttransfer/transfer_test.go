package nfttransfer_test

import (
	"math/big"
	"testing"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/nfttransfer"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/testvars"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../.env")
}

func TestNftTransfer(t *testing.T) {
	fromAddr := "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR"
	publicKey := "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7"
	privateKey := "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs"

	nonceReq, err := nonce.MakeNonceRequest(fromAddr)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	nonceInfo := nonce.NonceInfo{
		UseIndexer:    false,
		NonceReqs:     []*pb.NonceRequest{nonceReq},
		ValidatorAddr: testvars.TEST_GRPC_ADDRESS,
	}

	symbol := "$NFT+0000"
	itemID := big.NewInt(0)
	recipient := "CY7JXLDTwfqYUZsLJ58bHScFpK2gSgFLQy2CfAVGDBBt"
	baseFeeSymbol := "$ZRA+0000"
	baseFeeParts := "1000000000"

	// Contract fee fields (as applicable)
	var contractFeeID *string = nil
	var contractFeeAmountParts *big.Int = nil

	txn, err := nfttransfer.CreateNftTransfer(
		nonceInfo,
		symbol,
		itemID,
		recipient,
		publicKey,
		privateKey,
		baseFeeSymbol,
		baseFeeParts,
		contractFeeID,
		contractFeeAmountParts,
	)

	if err != nil {
		t.Fatalf("Error creating NFT transfer transaction: %v", err)
	}

	_, err = nfttransfer.SendNftTransferTxn(testvars.TEST_GRPC_ADDRESS+":50052", txn)
	if err != nil {
		t.Fatalf("Error sending NFT transfer transaction: %v", err)
	}

	if testvars.TEST_GRPC_ADDRESS == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

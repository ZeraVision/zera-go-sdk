package contract_test

import (
	"testing"
	"time"

	"github.com/ZeraVision/zera-go-sdk/contract"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTokenCreation(t *testing.T) {
	// fromAddr := "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR"
	// publicKey := "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7"
	// privateKey := "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs"
	fromAddr := "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR"
	publicKey := "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7"
	privateKey := "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs"

	nonceReq, err := nonce.MakeNonceRequest(fromAddr)
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

	baseFeeSymbol := "$ZRA+0000" // some ACE authorized token
	baseFeeParts := "1000000000" // a number of parts sufficient to cover the fee with your symbol

	// Construct token according to your preferences

	// Buld more complex objects
	// Governance
	govType := contract.GovernanceTypeHelper{
		Type: contract.Staggared,
		ProposalPeriod: &contract.ProposalPeriod{
			PeriodType:   contract.Months,
			VotingPeriod: 2,
		},
		Stages: []*contract.Stage{ // present only for staged governance type
			{
				PeriodType:  contract.Days,
				Length:      14,
				Break:       false,
				MaxApproved: 20,
			},
			{
				PeriodType: contract.Days,
				Length:     7,
				Break:      true,
				//* max approved not processed if break
			},
			{
				PeriodType:  contract.Months,
				Length:      1,
				Break:       false,
				MaxApproved: 10,
			},
			{
				PeriodType: contract.Days, // On remainder PeriodType doesnt matter
				Length:     0,             // 0 = remainder of days left in proposal period
				Break:      true,
			},
		},

		StartTimestamp: timestamppb.New(time.Now()),
	}

	allowedProposalContracts := []string{
		"$TEST+0000",
		"$TEST+0001",
	}

	allowedVoteContracts := []string{
		"$TEST+0000",
		"$TEST+0001",
	}

	fastQuorum := 50.1 // optional

	// Configure as required
	gov, err := contract.CreateGovernance(govType, 50.1, &fastQuorum, allowedProposalContracts, allowedVoteContracts, 1.234, nil, false)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	tokenData := &contract.TokenData{
		ContractVersion: 100000, //version 1.0.0
		ContractId:      "$TEST+0000",
		Symbol:          "TEST",
		Name:            "Test Token",
		Governance:      gov,
	}
	//

	txn, err := contract.CreateTokenTXN(nonceInfo, tokenData, publicKey, privateKey, baseFeeSymbol, baseFeeParts)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if grpcAddr == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

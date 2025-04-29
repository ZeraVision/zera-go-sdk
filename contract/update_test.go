package contract_test

import (
	"math/big"
	"testing"
	"time"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/contract"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/testvars"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestTokenUpdate demonstrates a complete token update process using this SDK.
// This covers optional and required configurations and can serve as a base example.
// This sample does not cover all possible configurations of a token but gives an idea of where to start configuring
func TestTokenUpdate(t *testing.T) {
	// Account and key setup
	fromAddr := "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR"
	publicKey := "r_A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7" // r key with update permissions required
	privateKey := "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs"

	// Nonce is required to avoid replay attacks
	nonceReq, err := nonce.MakeNonceRequest(fromAddr)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Configure nonce resolution through validator
	nonceInfo := nonce.NonceInfo{
		UseIndexer:    false,
		NonceReqs:     []*pb.NonceRequest{nonceReq},
		ValidatorAddr: testvars.TEST_GRPC_ADDRESS,
	}

	// Fee configuration
	baseFeeSymbol := "$ZRA+0000"
	baseFeeParts := "1000000000" // 1 ZRA specified for sample purposes

	// Which contract is being updated?
	contractID := "$TEST+0000"

	// Basic token identity and flags -- //! remember most of these are optional in an update!
	newNameValue := "New Test Token"
	var newName *string = &newNameValue
	var quashThreshold *uint32 = nil
	kycStatusValue := false
	var kycStatus *bool = &kycStatusValue
	immutableKycStatusValue := true
	var immutableKycStatus *bool = &immutableKycStatusValue

	// Core token data object
	tokenData := &contract.UpdateData{
		ContractVersion:    100001, // version 1.0.1 (ex, higher than current version) || standard is major.minor.patch (z.xx.yyy)
		ContractId:         contractID,
		Name:               newName,
		QuashThreshold:     quashThreshold,
		KycStatus:          kycStatus,
		ImmutableKycStatus: immutableKycStatus,
	}

	// Optional: Update Governance configuration for proposals and voting
	govType := contract.GovernanceTypeHelper{
		Type: contract.Staggared,
		ProposalPeriod: &contract.ProposalPeriod{
			PeriodType:   contract.Months,
			VotingPeriod: 2,
		},
		StartTimestamp: timestamppb.New(time.Now()),
	}

	allowedProposalContracts := []string{"$TEST+0000", "$TEST+0001"}
	allowedVoteContracts := []string{"$TEST+0000", "$TEST+0001"}
	fastQuorum := 50.1

	tokenData.Governance, err = contract.CreateGovernance(
		govType,
		50.1,
		&fastQuorum,
		allowedProposalContracts,
		allowedVoteContracts,
		1.234,
		nil,
		false,
	)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Optional: Update Compliance rules for address eligibility
	complianceConfig := [][]contract.ComplianceConfig{
		{
			{ContractID: "$TEST+0000", Level: 1},
			{ContractID: "$TEST+0001", Level: 2},
		},
		{
			{ContractID: "$TEST+0000", Level: 2},
		},
	}

	tokenData.TokenCompliance, err = contract.CreateCompliance(complianceConfig)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Optional: Update Fee model for contract usage
	feeConfig := contract.ContractFeeConfig{
		Type:      contract.FeeCurrencyEquivalent,
		Address:   fromAddr,
		Fee:       0.20,
		Burn:      39.5,
		Validator: 10.5,
		AllowedFeeInstrument: []string{
			"$TEST+0000",
			"TEST+0001",
		},
	}

	bigParts := big.NewInt(1_000_000_000) // customize to the number of parts in this token

	tokenData.ContractFees, err = contract.CreateContractFee(feeConfig, bigParts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Optional: Update restricted admin keys
	// Note: You can not update keys with a key authority higher or equal to your own (except self), they must stay the exact same config. Lower KeyWeight = Higher authority (0 max authority)
	rKey1 := "r_A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7"
	rKey2 := "r_B_c_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV"

	restrictedKeys := []contract.RestrictedConfig{
		{
			PublicKey:      helper.PublicKey{Single: &rKey1},
			UpdateContract: true,
			Mint:           true,
			Transfer:       true,
			Propose:        true,
			Vote:           true,
			CurEquiv:       true,
			KeyWeight:      0,
		},
		{
			PublicKey: helper.PublicKey{Single: &rKey2},
			Quash:     true,
			Transfer:  true,
			Propose:   true,
			Vote:      false,
			KeyWeight: 1,
		},
	}

	tokenData.RestrictedKeys, err = contract.CreateRestrictedKeys(restrictedKeys)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Optional: Scheduled operating cost percentages
	expenseConfig := []contract.ExpenseRatioConfig{
		{Month: 1, Day: 1, Percent: 1.234},
		{Month: 3, Day: 31, Percent: 1},
	}

	tokenData.ExpenseRatio, err = contract.CreateExpenseRatio(expenseConfig)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Optional: Arbitrary metadata for your own use
	customConfig := []contract.KeyValuePair{
		{Key: "cool key", Value: "cool data"},
		{Key: "crazy cool key", Value: "crazy cool data"},
	}
	tokenData.CustomParameters = contract.CreateCustomParameters(customConfig)

	// Final step: Build the token creation transaction
	txn, err := contract.UpdateContractTXN(nonceInfo, tokenData, publicKey, privateKey, baseFeeSymbol, baseFeeParts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = contract.SendUpdate(testvars.TEST_GRPC_ADDRESS, txn)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Output transaction hash for verification or tracking
	if testvars.TEST_GRPC_ADDRESS == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

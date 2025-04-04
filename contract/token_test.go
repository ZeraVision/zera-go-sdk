package contract_test

import (
	"testing"
	"time"

	"github.com/ZeraVision/zera-go-sdk/contract"
	"github.com/ZeraVision/zera-go-sdk/convert"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// TestTokenCreation demonstrates a complete token creation process using this SDK.
// This covers optional and required configurations and can serve as a base example.
// This sample does not cover all possible configurations of a token but gives an idea of where to start configuring
func TestTokenCreation(t *testing.T) {
	// Account and key setup
	fromAddr := "8ZfvifzSPMhhhivnH6NtaBXcmF3vsSaiB8KBULTetBcR"
	publicKey := "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7"
	privateKey := "2ap5CkCekErkqJ4UuSGAW1BmRRRNr8hXaebudv1j8TY6mJMSsbnniakorFGmetE4aegsyQAD8WX1N8Q2Y45YEBDs"

	// Nonce is required to avoid replay attacks
	nonceReq, err := nonce.MakeNonceRequest(fromAddr)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// For this demo we will get the nonce directly from a validator, but other services (ie ZV's indexer) can also provide this function
	grpcAddr := "routing.zera.vision" // and/or other validators

	// Configure nonce resolution through validator
	nonceInfo := nonce.NonceInfo{
		UseIndexer:    false,
		NonceReq:      nonceReq,
		ValidatorAddr: grpcAddr + ":50051",
	}

	// Fee configuration
	baseFeeSymbol := "$ZRA+0000"
	baseFeeParts := "1000000000" // 1 ZRA specified for sample purposes

	// Basic token identity and flags
	contractID := "$TEST+0000"
	symbol := "TEST"
	name := "Test Token"
	updateExpenseRatio := false
	updateContractFees := false
	var quashThreshold *uint32 = nil
	kycStatus := false
	immutableKycStatus := true
	var curEquivalentStart *float64 = nil

	// Core token data object
	tokenData := &contract.TokenData{
		ContractVersion:    100000, // version 1.0.0 || standard is major.minor.patch (z.xx.yyy)
		ContractId:         contractID,
		Symbol:             symbol,
		Name:               name,
		UpdateExpenseRatio: updateExpenseRatio,
		UpdateContractFees: updateContractFees,
		QuashThreshold:     quashThreshold,
		KycStatus:          kycStatus,
		ImmutableKycStatus: immutableKycStatus,
		CurEquivStart:      curEquivalentStart,
	}

	// Optional: Governance configuration for proposals and voting
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

	// Optional: Compliance rules for address eligibility
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

	// Required: Smallest divisible unit of the token
	var parts uint64 = 1000000000
	bigParts := convert.ToBigInt(parts)
	unitName := "smolpart"

	tokenData.Denomination, err = contract.CreateDenomination(bigParts, unitName)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Required: Total supply cap
	maxSupply := 123456.789
	tokenData.MaxSupply, err = contract.CreateMaxSupply(maxSupply, bigParts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Optional: Scheduled release of max supply over time
	releaseConfig := []contract.ReleaseScheduleConfig{
		{
			ReleaseDate: timestamppb.New(time.Now()),
			Amount:      maxSupply / 2,
		},
		{
			ReleaseDate: timestamppb.New(time.Now().AddDate(0, 6, 0)),
			Amount:      maxSupply / 2,
		},
	}

	tokenData.MaxSupplyRelease, err = contract.CreateMaxSupplyRelease(releaseConfig, bigParts, tokenData.MaxSupply)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Optional: Fee model for contract usage
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

	tokenData.ContractFees, err = contract.CreateContractFee(feeConfig, bigParts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Optional: Time-delayed restricted admin keys
	rKey1 := "r_A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7"
	rKey2 := "r_B_c_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV"

	restrictedKeys := []contract.RestrictedConfig{
		{
			PublicKey:      helper.PublicKey{Single: &rKey1},
			TimeDelay:      25,
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

	// Optional: Premint tokens to specific accounts
	premintConfig := []contract.PremintConfig{
		{Address: fromAddr, Amount: 123.456},
		{Address: "Hv3KUwrmR8C8XVSxuJFJrQqeDixeDnakUTkUUMZkFCUS", Amount: 789},
	}

	tokenData.Premint, err = contract.CreatePremint(premintConfig, bigParts)
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
	txn, err := contract.CreateTokenTXN(nonceInfo, tokenData, publicKey, privateKey, baseFeeSymbol, baseFeeParts)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, err = contract.SendInstrumentContract(grpcAddr, txn)

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Output transaction hash for verification or tracking
	if grpcAddr == "routing.zera.vision" {
		t.Logf("Transaction sent successfully: see (https://explorer.zera.vision/transactions/%s)", transcode.HexEncode(txn.Base.Hash))
	} else {
		t.Logf("Transaction sent successfully: %s", transcode.HexEncode(txn.Base.Hash))
	}
}

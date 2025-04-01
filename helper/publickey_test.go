package helper_test

import (
	"fmt"
	"testing"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"google.golang.org/protobuf/proto"
)

func TestGeneratePublicKey_SingleKey(t *testing.T) {
	expectedResult := "0a24415f635fd5c908ae57a79d4e820f61b20b618fdb782e87a46eaef580ec94c6815b82f90a"

	singleKey := "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7"
	publicKey := helper.PublicKey{
		Single: &singleKey,
	}

	result, err := helper.GeneratePublicKey(publicKey)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	compareResult(t, result, expectedResult)
}

func TestGeneratePublicKey_InheritenceKey(t *testing.T) {
	inheritenceKey := "$ZRA+0000"
	publicKey := helper.PublicKey{
		Inheritence: &inheritenceKey,
	}

	result, err := helper.GeneratePublicKey(publicKey)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	compareResult(t, result, "0a09245a52412b30303030")
}

func TestGeneratePublicKey_InvalidInheritenceKey(t *testing.T) {
	invalidKey := "InvalidKey"
	publicKey := helper.PublicKey{
		Inheritence: &invalidKey,
	}

	_, err := helper.GeneratePublicKey(publicKey)
	if err == nil {
		t.Fatal("Expected an error for invalid inheritence key, got none")
	}

	expectedError := "inheritence key must be in the format $LETTERS+0000 (with exactly 4 digits at the end)"
	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, got: %s", expectedError, err.Error())
	}
}

func TestGeneratePublicKey_MultiKey(t *testing.T) {
	multiKey := helper.MultiKeyHelper{
		MultiKey: []helper.MultiKey{
			{Class: 1, PublicKey: "A_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7"},
			{Class: 2, PublicKey: "B_8TZAaoUWbGvkxaWdWBXJ3mVHXVXLDJgtbeexkBzj5ySjpru7yZvfuKwGGHt2gtFpQfQCaRnBPU43bV"},
		},
		Pattern: [][]helper.MultiPatterns{
			{
				{Class: 1, Required: 1},
				{Class: 2, Required: 1},
			},
		},
		HashTokens: []helper.HashType{helper.BLAKE3, helper.SHA3_256},
	}

	publicKey := helper.PublicKey{
		Multi: &multiKey,
	}

	result, err := helper.GeneratePublicKey(publicKey)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	compareResult(t, result, "12710a22415fd5c908ae57a79d4e820f61b20b618fdb782e87a46eaef580ec94c6815b82f90a0a3b425f3e6467cc0cb5654788058c44db412b47d31f897ead16fc609b621e356c9455ed0fc6e1c59feab234f6385c331f84453bd6a06b6ad218545a001a080a02010212020101220163220161")
}

func TestGeneratePublicKey_SmartContractKey(t *testing.T) {
	smartContract := helper.SmartContractHelper{
		Name:     "TestContract",
		Instance: 1234,
	}

	publicKey := helper.PublicKey{
		SmartContract: &smartContract,
	}

	result, err := helper.GeneratePublicKey(publicKey)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	compareResult(t, result, "1a1473635f54657374436f6e74726163745f31323334")
}

func TestGeneratePublicKey_GovernanceKey(t *testing.T) {
	governanceKey := "$ZRA+0000"
	publicKey := helper.PublicKey{
		Governance: &governanceKey,
	}

	result, err := helper.GeneratePublicKey(publicKey)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	compareResult(t, result, "220d676f765f245a52412b30303030")
}

func TestGeneratePublicKey_MultipleKeysSet(t *testing.T) {
	singleKey := "A_c_FPXdqFTeqC3rHCaAAXmXbunb8C5BbRZEZNGjt23dAVo7"
	governanceKey := "$ZRA+0000"

	publicKey := helper.PublicKey{
		Single:     &singleKey,
		Governance: &governanceKey,
	}

	_, err := helper.GeneratePublicKey(publicKey)
	if err == nil {
		t.Fatal("Expected an error for multiple keys set, got none")
	}

	expectedError := "exactly one of the public key options must be set"
	if err.Error() != expectedError {
		t.Errorf("Expected error: %s, got: %s", expectedError, err.Error())
	}
}

// Converts protoobject to binary -> hex and compares it with expected result
func compareResult(t *testing.T, result *pb.PublicKey, expectedResult string) {
	// Serialize to binary
	binaryData, err := proto.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to serialize protobuf: %v", err)
	}

	compareResult := transcode.HexEncode(binaryData)

	fmt.Println(compareResult)

	if compareResult != expectedResult {
		t.Errorf("Expected result: %s, got: %s", expectedResult, compareResult)
	}
}

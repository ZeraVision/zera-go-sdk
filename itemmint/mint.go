package itemmint

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/ZeraVision/zera-go-sdk/wallet"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateItemMintTxn creates a new ItemizedMintTXN with the specified parameters.
// - votingWeight: optional voting weight (pointer to string)
// - contractFees: optional *pb.ItemContractFees
func CreateItemMintTxn(
	nonceInfo nonce.NonceInfo,
	contractId string,
	itemId *big.Int,
	recipient string,
	publicKeyBase58 string,
	privateKeyBase58 string,
	feeID string,
	feeAmountParts string,
	parameters []*pb.KeyValuePair,
	expiry *uint64,
	validFrom *uint64,
	votingWeight *big.Int,
	contractFees *pb.ItemContractFees,
) (*pb.ItemizedMintTXN, error) {
	// Step 1: Decode recipient address
	recipientBytes, err := transcode.Base58Decode(recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to decode recipient address: %v", err)
	}

	// Step 2: Decode public key
	_, _, pubKeyBytes, err := transcode.Base58DecodePublicKey(publicKeyBase58)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %v", err)
	}

	nonceArr, err := nonce.GetNonce(nonceInfo, 5)
	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}
	if len(nonceArr) != 1 {
		return nil, fmt.Errorf("expected exactly one nonce, got %d", len(nonceArr))
	}

	// Step 3: Create BaseTXN
	base := &pb.BaseTXN{
		PublicKey: &pb.PublicKey{},
		FeeId:     feeID,
		FeeAmount: feeAmountParts,
		Timestamp: timestamppb.New(time.Now().UTC()),
		Nonce:     nonceArr[0],
	}

	if strings.HasPrefix(publicKeyBase58, "gov_") {
		base.PublicKey.GovernanceAuth = pubKeyBytes
	} else if strings.HasPrefix(publicKeyBase58, "sc_") {
		base.PublicKey.SmartContractAuth = pubKeyBytes
	} else {
		base.PublicKey.Single = pubKeyBytes
	}

	// Step 4: Construct ItemizedMintTXN
	itemMintTxn := &pb.ItemizedMintTXN{
		Base:             base,
		ContractId:       contractId,
		ItemId:           itemId.Text(10),
		RecipientAddress: recipientBytes,
		Parameters:       parameters,
		ContractFees:     contractFees,
	}
	if expiry != nil {
		itemMintTxn.Expiry = expiry
	}
	if validFrom != nil {
		itemMintTxn.ValidFrom = validFrom
	}
	if votingWeight != nil {
		votingWeightStr := votingWeight.Text(10)
		itemMintTxn.VotingWeight = &votingWeightStr
	}

	// Step 5: Serialize transaction before signing
	byteDataNoSig, err := proto.Marshal(itemMintTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize item mint transaction: %v", err)
	}

	// Step 6: Verify and determine key
	if !strings.HasPrefix(publicKeyBase58, "r_") && !strings.HasPrefix(publicKeyBase58, "gov_") && !strings.HasPrefix(publicKeyBase58, "sc_") {
		return nil, fmt.Errorf("not possible to do restricted logic (requires r_, gov_, or sc_ key): %s", publicKeyBase58)
	}

	keyType, err := wallet.DetermineKeyType(publicKeyBase58)
	if err != nil {
		return nil, fmt.Errorf("failed to determine key type: %v", err)
	}

	signature, err := helper.Sign(privateKeyBase58, byteDataNoSig, keyType)
	if err != nil {
		return nil, fmt.Errorf("failed to sign item mint transaction: %v", err)
	}

	// Step 7: Assign signature to base
	itemMintTxn.Base.Signature = signature

	// Step 8: Serialize again with signature
	byteDataWithSig, err := proto.Marshal(itemMintTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed item mint transaction: %v", err)
	}

	// Step 9: Hash the signed transaction
	hash := transcode.SHA3256(byteDataWithSig)
	itemMintTxn.Base.Hash = hash

	return itemMintTxn, nil
}

// SendMintTXN submits a MintTXN to the network via gRPC
func SendItemMintTXN(grpcAddr string, txn *pb.ItemizedMintTXN) (*emptypb.Empty, error) {
	if !strings.Contains(grpcAddr, ":") {
		grpcAddr += ":50052"
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTXNServiceClient(conn)
	response, err := client.ItemMint(context.Background(), txn)
	if err != nil {
		return nil, fmt.Errorf("item mint transaction failed: %v", err)
	}

	return response, nil
}

// BuildItemContractFees builds a *pb.ItemContractFees from user-friendly types.
// fee: the fee amount as float64 (e.g., 1.23 for $1.23)
// feeAddressB58: the fee address in base58 encoding
// burnPercent: burn percent (0-100), e.g., 25.5 for 25.5%
// validatorPercent: validator percent (0-100), e.g., 10 for 10%
// allowedFeeInstruments: allowed fee instrument strings
func BuildItemContractFees(
	fee float64,
	feeAddressB58 string,
	burnPercent float64,
	validatorPercent float64,
	allowedFeeInstruments []string,
) (*pb.ItemContractFees, error) {
	const quintillion = 1e18

	// Convert fee to string in quintillion units
	feeInt := int64(math.Round(fee * quintillion))
	feeStr := fmt.Sprintf("%d", feeInt)

	// Decode fee address from base58
	feeAddrBytes, err := transcode.Base58Decode(feeAddressB58)
	if err != nil {
		return nil, fmt.Errorf("failed to decode fee address: %v", err)
	}

	// Convert burn and validator percent to string in quintillion units
	burnInt := int64(math.Round((burnPercent / 100.0) * quintillion))
	burnStr := fmt.Sprintf("%d", burnInt)

	validatorInt := int64(math.Round((validatorPercent / 100.0) * quintillion))
	validatorStr := fmt.Sprintf("%d", validatorInt)

	return &pb.ItemContractFees{
		Fee:                  feeStr,
		FeeAddress:           feeAddrBytes,
		Burn:                 burnStr,
		Validator:            validatorStr,
		AllowedFeeInstrument: allowedFeeInstruments,
	}, nil
}

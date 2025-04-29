package contract

import (
	"context"
	"fmt"
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

type UpdateData struct {
	ContractId         string                // Your existing contract ID
	ContractVersion    uint64                // [suggested] major (as needed) - minor (x2) - patch (x3) (ie 1.1.0 = 101000 = 1 | 01 | 000) //! version number must be greater than current version
	Name               *string               // Optionally change the name
	Memo               *string               // Memo in base
	Governance         *pb.Governance        // Governance configuration (ie staged, cycle, adaptive, staged, cycle, [empty])
	RestrictedKeys     []*pb.RestrictedKey   // Restricted keys for the contract (ie empty (not suitable for most cases), other specicial permissions)
	ContractFees       *pb.ContractFees      // Configuration to charge fees to users for *transfers* on top of base fees
	CustomParameters   []*pb.KeyValuePair    // Custom values contracts can put on chain (can be used for integrations / standards)
	ExpenseRatio       []*pb.ExpenseRatio    // Expense ratio configuration for contract (most contracts don't use this)
	TokenCompliance    []*pb.TokenCompliance // Token compliance configuration (most contracts don't use this)
	KycStatus          *bool                 // If true, this contract requires KYC (compliance status) to transact.
	ImmutableKycStatus *bool                 // If true, this contract cannot change its requirement of KYC status later.
	QuashThreshold     *uint32               // Number of restricted wallets needed to quash a transaction (most contracts don't use this)
}

func UpdateContractTXN(nonceInfo nonce.NonceInfo, data *UpdateData, publicKeyBase58 string, privateKeyBase58 string, feeID string, feeAmountParts string) (*pb.ContractUpdateTXN, error) {
	// Step 1: Decode public key
	prefix, _, pubKeyBytes, err := transcode.Base58DecodePublicKey(publicKeyBase58)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %v", err)
	}

	if !strings.HasPrefix(string(prefix), "r_") {
		return nil, fmt.Errorf("public key %s is not a restricted key", publicKeyBase58)
	}

	nonce, err := nonce.GetNonce(nonceInfo, 5)

	if err != nil {
		return nil, fmt.Errorf("failed to get nonce: %v", err)
	}

	if len(nonce) != 1 {
		return nil, fmt.Errorf("expected exactly one nonce, got %d", len(nonce))
	}

	// Step 2: Create BaseTXN
	base := &pb.BaseTXN{
		PublicKey: &pb.PublicKey{
			Single: pubKeyBytes,
		},
		FeeId:     feeID,
		FeeAmount: feeAmountParts,
		Timestamp: timestamppb.New(time.Now().UTC()),
		Nonce:     nonce[0],
		Memo:      data.Memo,
	}

	// Step 3: Construct Update TXN
	contractTxn := &pb.ContractUpdateTXN{
		Base:               base,
		ContractId:         data.ContractId,
		ContractVersion:    data.ContractVersion,
		Name:               data.Name,
		Governance:         data.Governance,
		RestrictedKeys:     data.RestrictedKeys,
		ContractFees:       data.ContractFees,
		CustomParameters:   data.CustomParameters,
		ExpenseRatio:       data.ExpenseRatio,
		TokenCompliance:    data.TokenCompliance,
		KycStatus:          data.KycStatus,
		ImmutableKycStatus: data.ImmutableKycStatus,
		QuashThreshold:     data.QuashThreshold,
	}

	// Step 4: Serialize transaction before signing
	byteDataNoSig, err := proto.Marshal(contractTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize contract transaction: %v", err)
	}

	keyType, err := wallet.DetermineKeyType(publicKeyBase58)
	if err != nil {
		return nil, fmt.Errorf("failed to determine key type: %v", err)
	}

	signature, err := helper.Sign(privateKeyBase58, byteDataNoSig, keyType)

	if err != nil {
		return nil, fmt.Errorf("failed to sign contract transaction: %v", err)
	}

	// Step 5: Assign signature to base
	contractTxn.Base.Signature = signature

	// Step 6: Serialize again with signature
	byteDataWithSig, err := proto.Marshal(contractTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed contract transaction: %v", err)
	}

	// Step 7: Hash the signed transaction
	hash := transcode.SHA3256(byteDataWithSig)
	contractTxn.Base.Hash = hash

	return contractTxn, nil
}

// SendUpdate submits an instrument contract to the network via gRPC
func SendUpdate(grpcAddr string, txn *pb.ContractUpdateTXN) (*emptypb.Empty, error) {

	if !strings.Contains(grpcAddr, ":") {
		grpcAddr += ":50052"
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTXNServiceClient(conn)
	response, err := client.ContractUpdate(context.Background(), txn)
	if err != nil {
		return nil, fmt.Errorf("token transaction failed: %v", err)
	}

	return response, nil
}

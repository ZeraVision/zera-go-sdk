package compliance

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ComplianceDetails struct {
	WalletAddr string
	Level      uint32
	Assign     bool // true to assign, false to revoke
	Expiry     *timestamp.Timestamp
}

func CreateComplianceTxn(nonceInfo nonce.NonceInfo, symbol string, details []ComplianceDetails, publicKeyBase58 string, privateKeyBase58 string, feeID string, feeAmountParts string) (*pb.ComplianceTXN, error) {
	// Step 1: Decode public key
	_, _, pubKeyBytes, err := transcode.Base58DecodePublicKey(publicKeyBase58)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %v", err)
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
	}

	// Step 3: Construct Compliance
	allowanceTxn := &pb.ComplianceTXN{
		Base:       base,
		ContractId: symbol,
		Compliance: []*pb.ComplianceAssign{},
	}

	for _, detail := range details {
		walletAddrByte, err := transcode.Base58Decode(detail.WalletAddr)
		if err != nil {
			return nil, fmt.Errorf("failed to decode wallet address: %v", err)
		}
		allowanceTxn.Compliance = append(allowanceTxn.Compliance, &pb.ComplianceAssign{
			RecipientAddress: walletAddrByte,
			ComplianceLevel:  detail.Level,
			AssignRevoke:     detail.Assign,
			Expiry:           detail.Expiry,
		})
	}

	// Step 4: Serialize transaction before signing
	byteDataNoSig, err := proto.Marshal(allowanceTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize allowance transaction: %v", err)
	}

	keyType, err := helper.DetermineKeyType(publicKeyBase58)

	if err != nil {
		return nil, fmt.Errorf("failed to determine key type: %v", err)
	}

	// Step 5: Sign the transaction
	signature, err := helper.Sign(privateKeyBase58, byteDataNoSig, keyType)

	if err != nil {
		return nil, fmt.Errorf("failed to sign compliance transaction: %v", err)
	}

	// Step 6: Assign signature to base
	allowanceTxn.Base.Signature = signature

	// Step 7: Serialize again with signature
	byteDataWithSig, err := proto.Marshal(allowanceTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed compliance transaction: %v", err)
	}

	// Step 8: Hash the signed transaction
	hash := transcode.SHA3256(byteDataWithSig)
	allowanceTxn.Base.Hash = hash

	return allowanceTxn, nil
}

func SendComplianceTxn(grpcAddr string, txn *pb.ComplianceTXN) (*emptypb.Empty, error) {
	if !strings.Contains(grpcAddr, ":") {
		grpcAddr += ":50052"
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTXNServiceClient(conn)
	response, err := client.Compliance(context.Background(), txn)
	if err != nil {
		return nil, fmt.Errorf("compliance transaction failed: %v", err)
	}

	return response, nil
}

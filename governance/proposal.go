package governance

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateProposalTxn(nonceInfo nonce.NonceInfo, symbol string, publicKeyBase58 string, privateKeyBase58 string, feeID string, feeAmountParts string, title, synopsis, body string, options []string, startTimestamp *timestamppb.Timestamp, endTimestamp *timestamppb.Timestamp, txns []*pb.GovernanceTXN) (*pb.GovernanceProposal, error) {
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

	// Step 3: Construct & Configure Proposal
	proposalTxn := &pb.GovernanceProposal{
		Base:           base,
		ContractId:     symbol,
		Title:          title,
		Synopsis:       synopsis,
		Body:           body,
		Options:        options,        // optional
		StartTimestamp: startTimestamp, // present for adaptive gov types
		EndTimestamp:   endTimestamp,   // present for adaptive gov types
		GovernanceTxn:  txns,
	}

	// Step 4: Serialize transaction before signing
	byteDataNoSig, err := proto.Marshal(proposalTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize proposal transaction: %v", err)
	}

	keyType, err := helper.DetermineKeyType(publicKeyBase58)

	if err != nil {
		return nil, fmt.Errorf("failed to determine key type: %v", err)
	}

	// Step 5: Sign the transaction
	signature, err := helper.Sign(privateKeyBase58, byteDataNoSig, keyType)

	if err != nil {
		return nil, fmt.Errorf("failed to sign proposal transaction: %v", err)
	}

	// Step 6: Assign signature to base
	proposalTxn.Base.Signature = signature

	// Step 7: Serialize again with signature
	byteDataWithSig, err := proto.Marshal(proposalTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed proposal transaction: %v", err)
	}

	// Step 8: Hash the signed transaction
	hash := transcode.SHA3256(byteDataWithSig)
	proposalTxn.Base.Hash = hash

	return proposalTxn, nil
}

// SendProposal submits a proposal to the network via gRPC
func SendProposal(grpcAddr string, txn *pb.GovernanceProposal) (*emptypb.Empty, error) {
	if !strings.Contains(grpcAddr, ":") {
		grpcAddr += ":50052"
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTXNServiceClient(conn)
	response, err := client.GovernProposal(context.Background(), txn)
	if err != nil {
		return nil, fmt.Errorf("proposal transaction failed: %v", err)
	}

	return response, nil
}

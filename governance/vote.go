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

// CreateVoteTxn creates a governance vote transaction by decoding the proposal ID and public key,
// retrieving the nonce, constructing the base transaction, and signing and hashing the transaction.
// It returns the constructed GovernanceVote protobuf object or an error if any step fails.
//
// Parameters:
// - nonceInfo: Information required to retrieve the nonce.
// - symbol: The contract symbol for the governance vote.
// - proposalID: The proposal ID in hexadecimal format.
// - publicKeyBase58: The voter's public key in Base58 format.
// - privateKeyBase58: The voter's private key in Base58 format.
// - feeID: The fee ID for the transaction.
// - feeAmountParts: The fee amount in parts.
// - support: A pointer to a boolean indicating whether the vote supports the proposal.
// - voteOption: A pointer to a uint32 specifying the vote option.
//
// Returns:
// - *pb.GovernanceVote: The constructed governance vote transaction.
// - error: An error if any step in the process fails.
func CreateVoteTxn(nonceInfo nonce.NonceInfo, symbol string, proposalID string, publicKeyBase58 string, privateKeyBase58 string, feeID string, feeAmountParts string, support *bool, voteOption *uint32) (*pb.GovernanceVote, error) {
	// Step 1: proposalID (from hex)
	proposalBytes, err := transcode.HexDecode(proposalID)
	if err != nil {
		return nil, fmt.Errorf("failed to decode proposalID: %v", err)
	}

	// Step 2: Decode public key
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

	// Step 3: Create BaseTXN
	base := &pb.BaseTXN{
		PublicKey: &pb.PublicKey{
			Single: pubKeyBytes,
		},
		FeeId:     feeID,
		FeeAmount: feeAmountParts,
		Timestamp: timestamppb.New(time.Now().UTC()),
		Nonce:     nonce[0],
	}

	// Step 4: Construct Vote
	voteTxn := &pb.GovernanceVote{
		Base:          base,
		ContractId:    symbol,
		ProposalId:    proposalBytes,
		Support:       support,
		SupportOption: voteOption,
	}

	// Step 5: Serialize transaction before signing
	byteDataNoSig, err := proto.Marshal(voteTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize vote transaction: %v", err)
	}

	keyType, err := helper.DetermineKeyType(publicKeyBase58)

	if err != nil {
		return nil, fmt.Errorf("failed to determine key type: %v", err)
	}

	// Step 6: Sign the transaction
	signature, err := helper.Sign(privateKeyBase58, byteDataNoSig, keyType)

	if err != nil {
		return nil, fmt.Errorf("failed to sign vote transaction: %v", err)
	}

	// Step 7: Assign signature to base
	voteTxn.Base.Signature = signature

	// Step 8: Serialize again with signature
	byteDataWithSig, err := proto.Marshal(voteTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed vote transaction: %v", err)
	}

	// Step 9: Hash the signed transaction
	hash := transcode.SHA3256(byteDataWithSig)
	voteTxn.Base.Hash = hash

	return voteTxn, nil
}

// SendVoteTxn submits a vote to the network via gRPC
func SendVoteTxn(grpcAddr string, txn *pb.GovernanceVote) (*emptypb.Empty, error) {
	if !strings.Contains(grpcAddr, ":") {
		grpcAddr += ":50052"
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTXNServiceClient(conn)
	response, err := client.GovernVote(context.Background(), txn)
	if err != nil {
		return nil, fmt.Errorf("vote transaction failed: %v", err)
	}

	return response, nil
}

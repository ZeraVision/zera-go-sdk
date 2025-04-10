package mint

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/sign"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateMintTxn creates a MintTXN protobuf message for minting new tokens to a specific address.
//
// Parameters:
// - useIndexer: true for data grab from indexer, false for validator
// - symbol: contract symbol to mint (example: $ZRA+0000)
// - amount: amount to mint, in full coins (as a string, not parts)
// - recipient: Base58-encoded address to receive the minted tokens
// - publicKeyBase58: Base58-encoded public key of the minting authority
// - privateKeyBase58: Base58-encoded private key corresponding to the above public key
// - feeID: fee id for the mint operation (example: $ZRA+0000)
// - feeAmountParts: fee amount in *parts* (example: 1000000000 = 1 ZRA)
//
// Returns:
// - *pb.MintTXN: the constructed and signed MintTXN
// - error: if any step in construction or signing fails
func CreateMintTxn(nonceInfo nonce.NonceInfo, symbol string, amount string, recipient string, publicKeyBase58 string, privateKeyBase58 string, feeID string, feeAmountParts string) (*pb.MintTXN, error) {
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

	nonce, err := nonce.GetNonce(nonceInfo)

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

	// Step 4: Construct MintTXN
	mintTxn := &pb.MintTXN{
		Base:             base,
		ContractId:       symbol,
		Amount:           amount,
		RecipientAddress: recipientBytes,
	}

	// Step 5: Serialize transaction before signing
	byteDataNoSig, err := proto.Marshal(mintTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize mint transaction: %v", err)
	}

	// Step 6: Verify and determine key
	// Check to ensure its a restricted key
	if !strings.HasPrefix(publicKeyBase58, "r_") {
		return nil, fmt.Errorf("public key is not a restricted key (r_): %s", publicKeyBase58)
	}

	// Find the substring between "r_" and the next "_"
	pubBeyParts := strings.SplitN(publicKeyBase58, "_", 3) // Split into at most 3 parts
	var keyLetter string
	if len(pubBeyParts) > 2 {
		keyLetter = pubBeyParts[1]
	} else {
		return nil, fmt.Errorf("invalid public key format: %s", publicKeyBase58)
	}

	// Ed25519
	var keyType helper.KeyType
	if keyLetter == "A" {
		keyType = helper.ED25519
	} else if keyLetter == "B" {
		keyType = helper.ED448
	} else {
		return nil, fmt.Errorf("unknown key type for public key: %s", publicKeyBase58)
	}

	signature, err := sign.Sign(privateKeyBase58, byteDataNoSig, keyType)

	if err != nil {
		return nil, fmt.Errorf("failed to sign mint transaction: %v", err)
	}

	// Step 7: Assign signature to base
	mintTxn.Base.Signature = signature

	// Step 8: Serialize again with signature
	byteDataWithSig, err := proto.Marshal(mintTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed mint transaction: %v", err)
	}

	// Step 9: Hash the signed transaction
	hash := transcode.SHA3256(byteDataWithSig)
	mintTxn.Base.Hash = hash

	return mintTxn, nil
}

// SendMintTXN submits a MintTXN to the network via gRPC
func SendMintTXN(grpcAddr string, txn *pb.MintTXN) (*emptypb.Empty, error) {
	if !strings.Contains(grpcAddr, ":") {
		grpcAddr += ":50052"
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTXNServiceClient(conn)
	response, err := client.Mint(context.Background(), txn)
	if err != nil {
		return nil, fmt.Errorf("mint transaction failed: %v", err)
	}

	return response, nil
}

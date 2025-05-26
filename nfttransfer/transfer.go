package nfttransfer

import (
	"context"
	"fmt"
	"math/big"
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

// CreateNftTransfer creates an NFT transfer transaction by decoding the item ID and recipient address,
// retrieving the nonce, constructing the base transaction, and signing and hashing the transaction.
// It returns the constructed NFTTXN protobuf object or an error if any step fails.
//
// Parameters:
// - nonceInfo: Information required to retrieve the nonce.
// - symbol: The contract symbol for the NFT.
// - itemID: The NFT item ID (string, as required by NFTTXN).
// - recipientBase58: The recipient's address in Base58 format.
// - publicKeyBase58: The sender's public key in Base58 format.
// - privateKeyBase58: The sender's private key in Base58 format.
// - feeID: The fee ID for the transaction.
// - feeAmountParts: The fee amount in parts.
// - contractFeeID: (if applicable) contract fee ID - fees associated with individual item
// - contractFeeAmountParts: (if applicable) contract fee amount in parts
//
// Returns:
// - *pb.NFTTXN: The constructed NFT transfer transaction.
// - error: An error if any step in the process fails.
func CreateNftTransfer(
	nonceInfo nonce.NonceInfo,
	symbol string,
	itemID *big.Int,
	recipientBase58 string,
	publicKeyBase58 string,
	privateKeyBase58 string,
	feeID string,
	feeAmountParts string,
	contractFeeID *string,
	contractFeeAmountParts *big.Int,
) (*pb.NFTTXN, error) {
	// Step 1: Decode recipient address
	recipientBytes, err := transcode.Base58Decode(recipientBase58)
	if err != nil {
		return nil, fmt.Errorf("failed to decode recipient address: %v", err)
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

	var contractFeeAmountPartsPtrStr *string
	if contractFeeAmountParts != nil {
		contractFeeAmountPartsStr := contractFeeAmountParts.Text(10)
		contractFeeAmountPartsPtrStr = &contractFeeAmountPartsStr
	}

	// Step 4: Construct NFTTXN
	nftTxn := &pb.NFTTXN{
		Base:              base,
		ContractId:        symbol,
		ItemId:            itemID.Text(10),
		RecipientAddress:  recipientBytes,
		ContractFeeId:     contractFeeID,
		ContractFeeAmount: contractFeeAmountPartsPtrStr,
	}

	// Step 5: Serialize transaction before signing
	byteDataNoSig, err := proto.Marshal(nftTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize NFT transfer transaction: %v", err)
	}

	keyType, err := helper.DetermineKeyType(publicKeyBase58)
	if err != nil {
		return nil, fmt.Errorf("failed to determine key type: %v", err)
	}

	// Step 6: Sign the transaction
	signature, err := helper.Sign(privateKeyBase58, byteDataNoSig, keyType)
	if err != nil {
		return nil, fmt.Errorf("failed to sign NFT transfer transaction: %v", err)
	}

	// Step 7: Assign signature to base
	nftTxn.Base.Signature = signature

	// Step 8: Serialize again with signature
	byteDataWithSig, err := proto.Marshal(nftTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed NFT transfer transaction: %v", err)
	}

	// Step 9: Hash the signed transaction
	hash := transcode.SHA3256(byteDataWithSig)
	nftTxn.Base.Hash = hash

	return nftTxn, nil
}

// SendNftTransferTxn submits an NFT transfer to the network via gRPC
func SendNftTransferTxn(grpcAddr string, txn *pb.NFTTXN) (*emptypb.Empty, error) {
	if !strings.Contains(grpcAddr, ":") {
		grpcAddr += ":50052"
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTXNServiceClient(conn)
	response, err := client.NFT(context.Background(), txn)
	if err != nil {
		return nil, fmt.Errorf("NFT transfer transaction failed: %v", err)
	}

	return response, nil
}

package currencyequivalent

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

type SelfData struct {
	Symbol string
	Rate   string
}

func CreateSelfCurrencyEquivalentTxn(nonceInfo nonce.NonceInfo, data []SelfData, publicKeyBase58 string, privateKeyBase58 string, feeID string, feeAmountParts string) (*pb.SelfCurrencyEquiv, error) {
	// Step 1: Decode public key
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

	var curEquiv []*pb.CurrencyEquiv

	for _, token := range data {
		curEquiv = append(curEquiv, &pb.CurrencyEquiv{
			ContractId: token.Symbol,
			Rate:       token.Rate,
		})
	}

	// Step 3: Construct ACE
	aceTxn := &pb.SelfCurrencyEquiv{
		Base:     base,
		CurEquiv: curEquiv,
	}

	// Step 4: Serialize transaction before signing
	byteDataNoSig, err := proto.Marshal(aceTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize ACE transaction: %v", err)
	}

	// Step 5: Verify and determine key
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

	signature, err := helper.Sign(privateKeyBase58, byteDataNoSig, keyType)

	if err != nil {
		return nil, fmt.Errorf("failed to sign mint transaction: %v", err)
	}

	// Step 6: Assign signature to base
	aceTxn.Base.Signature = signature

	// Step 7: Serialize again with signature
	byteDataWithSig, err := proto.Marshal(aceTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed mint transaction: %v", err)
	}

	// Step 8: Hash the signed transaction
	hash := transcode.SHA3256(byteDataWithSig)
	aceTxn.Base.Hash = hash

	return aceTxn, nil
}

func SendSelfCurrencyEquivalentTXN(grpcAddr string, txn *pb.SelfCurrencyEquiv) (*emptypb.Empty, error) {
	if !strings.Contains(grpcAddr, ":") {
		grpcAddr += ":50052"
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTXNServiceClient(conn)
	response, err := client.CurrencyEquiv(context.Background(), txn)
	if err != nil {
		return nil, fmt.Errorf("currency equivalent transaction failed: %v", err)
	}

	return response, nil
}

package allowance

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

type AllowanceDetails struct {
	Authorize          bool // true for approve, false for revoke
	WalletAddr         string
	CurrencyEquivalent *float64 // actual currency equivalent value, it 1.234, scaled within function (one of this OR amount present)
	Amount             *big.Int // parts of a token (one of this OR currency equivalent present)
	PeriodMonths       *uint32  // Number of months (one of this OR seconds present)
	PeriodSeconds      *uint32  // Number of seconds (one of this OR months present)
	StartTime          int64    // unix of starttime
}

func CreateAllowanceTxn(nonceInfo nonce.NonceInfo, symbol string, details AllowanceDetails, publicKeyBase58 string, privateKeyBase58 string, feeID string, feeAmountParts string) (*pb.AllowanceTXN, error) {
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

	// Decode b58 addr
	walletAddrByte, err := transcode.Base58Decode(details.WalletAddr)

	if err != nil {
		return nil, fmt.Errorf("failed to decode wallet address: %v", err)
	}

	var currencyEquivScaledStr *string

	// Take curEquivAmount to 1e18
	if details.CurrencyEquivalent != nil {
		currencyEquivScaled := new(big.Float).Mul(big.NewFloat(*details.CurrencyEquivalent), big.NewFloat(1e18)).String()
		currencyEquivScaledStr = &currencyEquivScaled
	}

	var amountParts *string
	if details.Amount != nil {
		tmp := details.Amount.String()
		amountParts = &tmp
	}

	if details.CurrencyEquivalent != nil && details.Amount != nil {
		return nil, fmt.Errorf("only one of CurrencyEquivalent or Amount should be provided")
	}

	if details.PeriodMonths != nil && details.PeriodSeconds != nil {
		return nil, fmt.Errorf("only one of PeriodMonths or PeriodSeconds should be provided")
	}

	var startTime *timestamppb.Timestamp
	if details.StartTime < 0 {
		details.StartTime = 0
	}

	startTime = timestamppb.New(time.Unix(details.StartTime, 0))

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

	// Step 3: Construct Vote
	allowanceTxn := &pb.AllowanceTXN{
		Base:       base,
		ContractId: symbol,

		Authorize:                 details.Authorize,
		WalletAddress:             walletAddrByte,
		AllowedCurrencyEquivelent: currencyEquivScaledStr,
		AllowedAmount:             amountParts,
		PeriodMonths:              details.PeriodMonths,
		PeriodSeconds:             details.PeriodSeconds,
		StartTime:                 startTime,
	}

	// nil unneeded params
	if !allowanceTxn.Authorize {
		allowanceTxn.AllowedCurrencyEquivelent = nil
		allowanceTxn.AllowedAmount = nil
		allowanceTxn.PeriodMonths = nil
		allowanceTxn.PeriodSeconds = nil
		allowanceTxn.StartTime = nil
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
		return nil, fmt.Errorf("failed to sign allowance transaction: %v", err)
	}

	// Step 6: Assign signature to base
	allowanceTxn.Base.Signature = signature

	// Step 7: Serialize again with signature
	byteDataWithSig, err := proto.Marshal(allowanceTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize signed allowance transaction: %v", err)
	}

	// Step 8: Hash the signed transaction
	hash := transcode.SHA3256(byteDataWithSig)
	allowanceTxn.Base.Hash = hash

	return allowanceTxn, nil
}

func SendAllowanceTxn(grpcAddr string, txn *pb.AllowanceTXN) (*emptypb.Empty, error) {
	if !strings.Contains(grpcAddr, ":") {
		grpcAddr += ":50052"
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTXNServiceClient(conn)
	response, err := client.Allowance(context.Background(), txn)
	if err != nil {
		return nil, fmt.Errorf("allowance transaction failed: %v", err)
	}

	return response, nil
}

package contract

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

type TokenData struct {
	ContractVersion    uint64                 // [suggested] major (as needed) - minor (x2) - patch (x3) (ie 1.1.0 = 101000 = 1 | 01 | 000)
	ContractId         string                 // Contract ID (ie $ZRA+0000 - must be unique to network)
	Symbol             string                 // Symbol (ie ZRA) without any contractID identifier (ie $ZRA+0000)
	Name               string                 // Name of the contract (ie ZERA)
	Governance         *pb.Governance         // Governance configuration (ie staged, cycle, adaptive, staged, cycle, [empty])
	RestrictedKeys     []*pb.RestrictedKey    // Restricted keys for the contract (ie empty (not suitable for most cases), other specicial permissions)
	Denomination       *pb.CoinDenomination   // How many 'parts' per coin there are (ie ZERA has 1_000_000_000 parts per coin)
	MaxSupply          string                 // Maximum supply of the contract in denomination units (not full coins)
	MaxSupplyRelease   []*pb.MaxSupplyRelease // An unlocking schedule
	Premint            []*pb.PreMintWallet    // Wallets to immediately premint supply to
	CustomParameters   []*pb.KeyValuePair     // Custom values contracts can put on chain (can be used for integrations / standards)
	ExpenseRatio       []*pb.ExpenseRatio     // Expense ratio configuration for contract (most contracts don't use this)
	UpdateExpenseRatio bool                   // Whether the contract can update its expense ratio later
	ContractFees       *pb.ContractFees       // Configuration to charge fees to users for *transfers* on top of base fees
	UpdateContractFees bool                   // Whether the contract can update its fees later
	QuashThreshold     *uint32                // Number of restricted wallets needed to quash a transaction (most contracts don't use this)
	TokenCompliance    []*pb.TokenCompliance  // Token compliance configuration (most contracts don't use this)
	KycStatus          bool                   // If true, this contract requires KYC (compliance status) to transact.
	ImmutableKycStatus bool                   // If true, this contract cannot change its requirement of KYC status later.
	CurEquivStart      *float64               // A starter version of "SelfCurrencyEquiv" that can set initial on chain rate, pass in as float64
}

func CreateTokenTXN(nonceInfo nonce.NonceInfo, data *TokenData, publicKeyBase58 string, privateKeyBase58 string, feeID string, feeAmountParts string) (*pb.InstrumentContract, error) {
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

	var startCurequivStr *string

	if data.CurEquivStart != nil {
		if *data.CurEquivStart < 0 {
			return nil, fmt.Errorf("CurEquivStart must be greater than 0")
		}

		// Convert to format network expects (1e18 scale)
		startCurequiv := new(big.Float).Mul(big.NewFloat(*data.CurEquivStart), big.NewFloat(1e18))

		startCurequivStrValue := startCurequiv.Text('f', 0)
		startCurequivStr = &startCurequivStrValue
	}

	// Step 3: Construct Token Contract
	contractTxn := &pb.InstrumentContract{
		Base:               base,
		Type:               pb.CONTRACT_TYPE_TOKEN,
		ContractVersion:    data.ContractVersion,
		ContractId:         data.ContractId,
		Symbol:             data.Symbol,
		Name:               data.Name,
		Governance:         data.Governance,
		RestrictedKeys:     data.RestrictedKeys,
		CoinDenomination:   data.Denomination,
		MaxSupply:          &data.MaxSupply,
		MaxSupplyRelease:   data.MaxSupplyRelease,
		PremintWallets:     data.Premint,
		CustomParameters:   data.CustomParameters,
		ExpenseRatio:       data.ExpenseRatio,
		UpdateExpenseRatio: data.UpdateExpenseRatio,
		ContractFees:       data.ContractFees,
		UpdateContractFees: data.UpdateContractFees,
		QuashThreshold:     data.QuashThreshold,
		TokenCompliance:    data.TokenCompliance,
		KycStatus:          data.KycStatus,
		ImmutableKycStatus: data.ImmutableKycStatus,
		CurEquivStart:      startCurequivStr,
	}

	// Step 4: Serialize transaction before signing
	byteDataNoSig, err := proto.Marshal(contractTxn)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize contract transaction: %v", err)
	}

	// Ed25519
	var keyType helper.KeyType
	if strings.HasPrefix(publicKeyBase58, "A") {
		keyType = helper.ED25519
	} else if strings.HasPrefix(publicKeyBase58, "B") {
		keyType = helper.ED448
	} else {
		return nil, fmt.Errorf("unknown key type for public key: %s", publicKeyBase58)
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

// SendInstrumentContract submits an instrument contract to the network via gRPC
func SendInstrumentContract(grpcAddr string, txn *pb.InstrumentContract) (*emptypb.Empty, error) {

	if !strings.Contains(grpcAddr, ":") {
		grpcAddr += ":50052"
	}

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewTXNServiceClient(conn)
	response, err := client.Contract(context.Background(), txn)
	if err != nil {
		return nil, fmt.Errorf("token transaction failed: %v", err)
	}

	return response, nil
}

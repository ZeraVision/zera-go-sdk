package transfer

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/nonce"
	"github.com/ZeraVision/zera-go-sdk/parts"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Inputs struct {
	B58Address         string
	KeyType            helper.KeyType
	PublicKey          string     // Base 58 encoded
	PrivateKey         string     // Base 58 encoded
	Amount             *big.Float // full coins (not parts)
	FeePercent         float32    // 0-100 max 6 digits of precision
	ContractFeePercent *float32   // 0-100 max 6 digits of precision
}

func CreateCoinTxn(nonceInfo nonce.NonceInfo, partsInfo parts.PartsInfo, inputs []Inputs, outputs map[string]*big.Float, baseFeeID, baseFeeAmountParts string, contractFeeID, contractFeeAmountParts *string) (*pb.CoinTXN, error) {

	parts, err := parts.GetParts(partsInfo)

	if err != nil {
		return nil, fmt.Errorf("could not get parts: %v", err)
	}

	// Step 1: Process Inputs
	inputTransfers, auth, keys, totalInput, err := processInputs(nonceInfo, inputs, parts)
	if err != nil {
		return nil, err
	}

	// Step 2: Process Outputs
	outputTransfers, totalOutput, err := processOutputs(outputs, parts)
	if err != nil {
		return nil, err
	}

	// Check to see if inputs and outputs match
	if totalInput.Cmp(totalOutput) != 0 {
		return nil, fmt.Errorf("total input does not equal total output: %s != %s", totalInput.String(), totalOutput.String())
	}

	// Step 3: Build Transfer Authentication
	transferAuth := buildTransferAuthentication(auth)

	// Step 4: Build Transaction Base
	txnBase := buildTransactionBase(baseFeeID, baseFeeAmountParts)

	// Step 5: Assemble Transaction
	txn := &pb.CoinTXN{
		Auth:              transferAuth,
		Base:              txnBase,
		ContractId:        partsInfo.Symbol,
		InputTransfers:    inputTransfers,
		OutputTransfers:   outputTransfers,
		ContractFeeId:     contractFeeID,
		ContractFeeAmount: contractFeeAmountParts,
	}

	// Step 6: Serialize and Sign Transaction
	txn, err = signTransaction(txn, keys)
	if err != nil {
		return nil, err
	}

	// Step 7: Marshal the vote with the signature
	byteDataWithSig, err := proto.Marshal(txn)
	if err != nil {
		return nil, fmt.Errorf("error while serializing txn: %v", err)
	}

	// Step 8: Hash the serialized data with the signature
	hash := transcode.SHA3256(byteDataWithSig)

	// Step 9: Add the hash
	txn.Base.Hash = hash

	return txn, nil
}

type keyTracking struct {
	KeyType    helper.KeyType
	PrivateKey string
}

type authTracking struct {
	PublicKeyBytes []byte
	Signature      []byte
	Nonce          uint64
}

// Helper Function: Process Inputs
// useIndexer - true if using indexer, false if using validator
// nonceEndpoint - url / addr of indexer or validator
func processInputs(nonceInfo nonce.NonceInfo, inputs []Inputs, parts *big.Int) ([]*pb.InputTransfers, []authTracking, map[string]keyTracking, *big.Float, error) {
	var inputTransfers []*pb.InputTransfers
	var auth []authTracking
	keys := map[string]keyTracking{}
	totalInput := big.NewFloat(0)
	index := uint64(0)

	for i, input := range inputs {
		_, _, pubKeyByte, err := transcode.Base58DecodePublicKey(input.PublicKey)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("could not decode public key: %v", err)
		}

		amountPartsBigF := new(big.Float).Mul(input.Amount, big.NewFloat(float64(parts.Int64())))

		inputTransfers = append(inputTransfers, &pb.InputTransfers{
			Index:      index,
			Amount:     amountPartsBigF.String(),
			FeePercent: uint32(input.FeePercent * 1_000_000),
		})

		nonce, err := nonce.GetNonce(nonceInfo)

		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("could not get nonce: %v", err)
		}

		auth = append(auth, authTracking{
			PublicKeyBytes: pubKeyByte,
			Signature:      nil,
			Nonce:          nonce[i],
		})

		keys[transcode.Base58Encode(pubKeyByte)] = keyTracking{
			KeyType:    input.KeyType,
			PrivateKey: input.PrivateKey,
		}

		totalInput.Add(totalInput, amountPartsBigF)
		index++
	}

	return inputTransfers, auth, keys, totalInput, nil
}

// Helper Function: Process Outputs
func processOutputs(outputs map[string]*big.Float, parts *big.Int) ([]*pb.OutputTransfers, *big.Float, error) {
	var outputsTransfers []*pb.OutputTransfers
	totalOutput := big.NewFloat(0)

	for address, amount := range outputs {
		decodedAddr, err := transcode.Base58Decode(address)
		if err != nil {
			return nil, nil, fmt.Errorf("could not decode address: %v", err)
		}

		bigFParts := new(big.Float).Mul(amount, new(big.Float).SetInt(parts))

		outputsTransfers = append(outputsTransfers, &pb.OutputTransfers{
			WalletAddress: decodedAddr,
			Amount:        bigFParts.String(),
		})

		totalOutput.Add(totalOutput, bigFParts)
	}

	return outputsTransfers, totalOutput, nil
}

// Helper Function: Build Transfer Authentication
func buildTransferAuthentication(auth []authTracking) *pb.TransferAuthentication {
	transferAuth := &pb.TransferAuthentication{}
	for _, a := range auth {
		transferAuth.PublicKey = append(transferAuth.PublicKey, &pb.PublicKey{Single: a.PublicKeyBytes})
		transferAuth.Nonce = append(transferAuth.Nonce, a.Nonce)
	}
	return transferAuth
}

// Helper Function: Build Transaction Base
func buildTransactionBase(feeID, feeAmountParts string) *pb.BaseTXN {
	return &pb.BaseTXN{
		Timestamp: timestamppb.Now(),
		FeeAmount: feeAmountParts,
		FeeId:     feeID,
	}
}

// Helper Function: Sign Transaction
func signTransaction(txn *pb.CoinTXN, keys map[string]keyTracking) (*pb.CoinTXN, error) {
	txnBytes, err := proto.Marshal(txn)
	if err != nil {
		return nil, fmt.Errorf("could not marshal transaction: %v", err)
	}

	for _, auth := range txn.Auth.PublicKey {
		if key, ok := keys[transcode.Base58Encode(auth.Single)]; ok {
			signature, err := helper.Sign(key.PrivateKey, txnBytes, key.KeyType)
			if err != nil {
				return nil, fmt.Errorf("could not sign transaction: %v", err)
			}
			txn.Auth.Signature = append(txn.Auth.Signature, signature)
		} else {
			return nil, fmt.Errorf("could not find private key for public key: %s", transcode.Base58Encode(auth.Single))
		}
	}
	return txn, nil
}

func SendCoinTXN(grpcAddr string, txn *pb.CoinTXN) (*emptypb.Empty, error) {
	if !strings.Contains(grpcAddr, ":") {
		grpcAddr += ":50052"
	}

	// Create a gRPC connection to the server
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// Create a new instance of ValidatorNetworkClient
	client := helper.NewNetworkClient(conn)

	response, err := client.Coin(context.Background(), txn)

	if err != nil {
		return nil, err
	}

	return response, nil
}

func CreateNonceRequests(inputs []Inputs) ([]*pb.NonceRequest, error) {
	var nonceReqs []*pb.NonceRequest

	for _, input := range inputs {
		nonceReq, err := nonce.MakeNonceRequest(input.B58Address)
		if err != nil {
			return nil, fmt.Errorf("error creating nonce request for address %s: %v", input.B58Address, err)
		}
		nonceReqs = append(nonceReqs, nonceReq)
	}

	return nonceReqs, nil
}

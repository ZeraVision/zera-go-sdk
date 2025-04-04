package contract

import (
	"fmt"
	"math/big"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/convert"
	"github.com/ZeraVision/zera-go-sdk/transcode"
)

type FeeType int16

const (
	FeeFixed              FeeType = 0 // an amount of parts to be paid
	FeeCurrencyEquivalent FeeType = 1 // rate based on self/ace oracles
	FeePercentage         FeeType = 2 // percentage of transaction amount
)

type ContractFeeConfig struct {
	Type                 FeeType  // the type of fee
	Address              string   // where the fees go to
	Fee                  float64  // depends on type: fixed (number of tokens) | current equivalent $xx.xx | percent (0-100) || remainder not given to burn and validator goes here
	Burn                 float64  // Percentage burned (0-100)
	Validator            float64  // Percentage to validator (0-100)
	AllowedFeeInstrument []string // contractID of the contracts allowed to pay the fee instrument. //! If allowed fee instrument is not self, a calculation is done based on self/auth currency equivalent values
}

func CreateContractFee(config ContractFeeConfig, parts *big.Int) (*pb.ContractFees, error) {
	feeAddr, err := transcode.Base58Decode(config.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to decode fee address: %v", err)
	}

	var feeString, burnString, validatorString string

	if config.Type == FeeFixed {
		fee := new(big.Float).Mul(big.NewFloat(config.Fee), convert.ToBigFloat(parts))
		feeInt := new(big.Int)
		feeInt, _ = fee.Int(feeInt)

		feeString = feeInt.String()
	} else if config.Type == FeeCurrencyEquivalent || config.Type == FeePercentage {
		// Scaled to 1e18 for network (0-100 scale * 1e16)
		fee := new(big.Float).Mul(big.NewFloat(config.Fee), convert.ToBigFloat(1e16))
		feeInt := new(big.Int)
		feeInt, _ = fee.Int(feeInt)

		feeString = feeInt.String()
	}

	// Burn percent
	burn := new(big.Float).Mul(big.NewFloat(config.Burn), convert.ToBigFloat(1e16))
	burnInt := new(big.Int)
	burnInt, _ = burn.Int(burnInt)

	burnString = burnInt.String()

	// Validator Percent
	validator := new(big.Float).Mul(big.NewFloat(config.Validator), convert.ToBigFloat(1e16))
	validatorInt := new(big.Int)
	validatorInt, _ = validator.Int(validatorInt)

	validatorString = validatorInt.String()

	return &pb.ContractFees{
		ContractFeeType:      pb.CONTRACT_FEE_TYPE(config.Type),
		Fee:                  feeString,
		Burn:                 burnString,
		Validator:            validatorString,
		FeeAddress:           feeAddr,
		AllowedFeeInstrument: config.AllowedFeeInstrument,
	}, nil
}

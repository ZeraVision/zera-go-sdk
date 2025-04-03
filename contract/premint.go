package contract

import (
	"fmt"
	"math/big"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/convert"
	"github.com/ZeraVision/zera-go-sdk/transcode"
)

type PremintConfig struct {
	Address string  // the address to premint in base58
	Amount  float64 // the amount as a whole number (parts calculation done within helper function)
}

func CreatePremint(premints []PremintConfig, parts *big.Int) ([]*pb.PreMintWallet, error) {

	var premintResult []*pb.PreMintWallet

	for _, premint := range premints {
		addr, err := transcode.Base58Decode(premint.Address)

		if err != nil {
			return nil, err
		}

		amountF := new(big.Float).Mul(big.NewFloat(premint.Amount), convert.ToBigFloat(parts))
		amount := new(big.Int)
		amount, accuracy := amountF.Int(amount)

		if accuracy != big.Exact {
			// Handle truncation or rounding
			return nil, fmt.Errorf("premint amount %v has more precision than specified in the contract (ie specified 0.999999999999999999 when maximum precision is only 0.9999999)", amountF)
		}

		premintResult = append(premintResult, &pb.PreMintWallet{
			Address: addr,
			Amount:  amount.String(),
		})
	}

	return premintResult, nil
}

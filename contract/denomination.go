package contract

import (
	"fmt"
	"math/big"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
)

func CreateDenomination(parts *big.Int, partName string) (*pb.CoinDenomination, error) {
	if parts.Cmp(big.NewInt(1)) < 0 {
		return nil, fmt.Errorf("denomination parts must be greater than 0")
	}

	// Create the denomination object
	return &pb.CoinDenomination{
		DenominationName: partName,
		Amount:           parts.String(),
	}, nil
}

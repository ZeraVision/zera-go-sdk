package contract

import (
	"fmt"
	"math/big"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/convert"
	"github.com/ZeraVision/zn-wallet-manager/functions"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CreateMaxSupply(maxSupply float64, parts *big.Int) (string, error) {
	// Convert into parts
	maxSupplyF := new(big.Float).Mul(big.NewFloat(maxSupply), convert.ToBigFloat(parts))
	amount := new(big.Int)
	amount, accuracy := maxSupplyF.Int(amount)
	if accuracy != big.Exact {
		return "", fmt.Errorf("max supply %v has more precision than specified in the contract (ie specified 0.999999999999999999 when maximum precision is only 0.9999999)", maxSupplyF)
	}

	return amount.String(), nil
}

type ReleaseScheduleConfig struct {
	ReleaseDate *timestamppb.Timestamp // the date of the release (in UTC)
	Amount      float64                // the amount to release
}

func CreateMaxSupplyRelease(releaseConfig []ReleaseScheduleConfig, parts *big.Int, maxSupply string) ([]*pb.MaxSupplyRelease, error) {
	var maxSupplyRelease []*pb.MaxSupplyRelease

	totalRelease := big.NewInt(0)

	for _, release := range releaseConfig {
		releaseParts := new(big.Float).Mul(big.NewFloat(release.Amount), convert.ToBigFloat(parts))
		releaseAmount := new(big.Int)
		releaseAmount, accuracy := releaseParts.Int(releaseAmount)
		if accuracy != big.Exact {
			return nil, fmt.Errorf("max supply release amount %v has more precision than specified in the contract (ie specified 0.999999999999999999 when maximum precision is only 0.9999999)", releaseParts)
		}

		maxSupplyRelease = append(maxSupplyRelease, &pb.MaxSupplyRelease{
			ReleaseDate: release.ReleaseDate,
			Amount:      releaseAmount.String(),
		})

		totalRelease = totalRelease.Add(totalRelease, releaseAmount)
	}

	if totalRelease.Cmp(functions.ToBigInt(maxSupply)) != 0 {
		return nil, fmt.Errorf("total release amount %v does not match max supply %v", totalRelease, maxSupply)
	}

	return maxSupplyRelease, nil
}

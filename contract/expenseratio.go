package contract

import (
	"fmt"
	"math"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
)

type ExpenseRatioConfig struct {
	Month   uint32  // month number (1 January - 12 December)
	Day     uint32  // day of month
	Percent float64 // 0 - 100
}

func CreateExpenseRatio(config []ExpenseRatioConfig) ([]*pb.ExpenseRatio, error) {

	var expenseRatio []*pb.ExpenseRatio

	for _, ratio := range config {
		// Validate that Percent has no more than 4 decimal places
		if math.Round(ratio.Percent*10000)/10000 != ratio.Percent {
			return nil, fmt.Errorf("percent value %.8f can not be more than 4 decimal places", ratio.Percent)
		}

		// Scale the percent value
		scaledPercent := uint32(ratio.Percent * 10000)

		// Append to the result
		expenseRatio = append(expenseRatio, &pb.ExpenseRatio{
			Day:     ratio.Day,
			Month:   ratio.Month,
			Percent: scaledPercent,
		})
	}

	return expenseRatio, nil
}

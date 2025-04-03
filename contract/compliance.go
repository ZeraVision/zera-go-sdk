package contract

import (
	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
)

type ComplianceConfig struct {
	ContractID string // what contract is the compliance certificate issued from?
	Level      uint32 // what rating is it
}

func CreateCompliance(config [][]ComplianceConfig) ([]*pb.TokenCompliance, error) {
	var complianceArr []*pb.TokenCompliance
	for _, outer := range config {
		var innerArr []*pb.Compliance
		for _, inner := range outer {
			innerArr = append(innerArr, &pb.Compliance{
				ContractId:      inner.ContractID,
				ComplianceLevel: inner.Level,
			})
		}
		complianceArr = append(complianceArr, &pb.TokenCompliance{
			Compliance: innerArr,
		})
	}

	return complianceArr, nil
}

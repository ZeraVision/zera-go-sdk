package contract

import (
	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
)

type KeyValuePair struct {
	Key   string
	Value string
}

func CreateCustomParameters(params []KeyValuePair) []*pb.KeyValuePair {
	var customParams []*pb.KeyValuePair

	for _, kvp := range params {
		customParams = append(customParams, &pb.KeyValuePair{
			Key:   kvp.Key,
			Value: kvp.Value,
		})
	}

	return customParams
}

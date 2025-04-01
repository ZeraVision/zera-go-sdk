package contract

import (
	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/helper"
)

// Only single keys supported in this function
func CreateRestrictedKey(publicKey helper.PublicKey) (*pb.RestrictedKey, error) {
	pubKey, err := helper.GeneratePublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	rKey := &pb.RestrictedKey{
		PublicKey: pubKey,
	}

	return rKey, nil
}

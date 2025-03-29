package helper

import (
	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"google.golang.org/grpc"
)

// struct for client implementation of grpcs
type NetworkClient struct {
	client pb.TXNServiceClient
}

// constructor for client implementation of grpcs
func NewNetworkClient(conn *grpc.ClientConn) pb.TXNServiceClient {
	return pb.NewTXNServiceClient(conn)
}

func NewValidatorNetworkClient(conn *grpc.ClientConn) pb.ValidatorServiceClient {
	return pb.NewValidatorServiceClient(conn)
}

package helper

import (
	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"google.golang.org/grpc"
)

func NewNetworkClient(conn *grpc.ClientConn) pb.TXNServiceClient {
	return pb.NewTXNServiceClient(conn)
}

func NewValidatorNetworkClient(conn *grpc.ClientConn) pb.ValidatorServiceClient {
	return pb.NewValidatorServiceClient(conn)
}

func NewValidatorNetworkApiClient(conn *grpc.ClientConn) pb.APIServiceClient {
	return pb.NewAPIServiceClient(conn)
}

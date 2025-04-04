package nonce

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NonceInfo struct {
	UseIndexer    bool             // use the ZV indexer if true, false use the validator gRPC api
	NonceReq      *pb.NonceRequest // required when useIndexer false
	ValidatorAddr string           // required when useIndexer false
	Address       string           // required when useIndexer true
	IndexerURL    string           // required when useIndexer true
	Authorization string           // required when useIndexer true, Api-Key or Bearer
	Override      uint64           // optional, if set, this nonce will be used instead of the one from the indexer or validator
}

// GetNonce retrieves the nonce either from the Indexer HTTP API or the Validator gRPC service.
// If useIndexer is true, indexerURL and apiKey must be provided. Uses ZV indexer (higher reliability, multiple geo locations (for lower global latency))
// If useIndexer is false, nonceReq and validatorAddr must be provided. Uses direct validator gRPC (lower reliability)
func GetNonce(info NonceInfo) (uint64, error) {

	if info.Override != 0 {
		return info.Override, nil
	}

	if info.UseIndexer {
		if info.IndexerURL == "" || info.Authorization == "" || info.Address == "" {
			return 0, fmt.Errorf("indexerURL, authorization, and address are required when useIndexer is true")
		}

		tr := &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
		client := &http.Client{Transport: tr}

		url := fmt.Sprintf("%s/store?requestType=getNextNonce&address=%s", info.IndexerURL, info.Address)

		req, err := http.NewRequest("POST", url, nil)
		if err != nil {
			return 0, fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Add("Target", "indexer")

		// Bearer
		if strings.Contains(info.Authorization, ".") {
			req.Header.Add("Authorization", "Bearer "+info.Authorization)
		} else { // Api Key
			req.Header.Add("Authorization", "Api-Key "+info.Authorization)
		}

		resp, err := client.Do(req)
		if err != nil {
			return 0, fmt.Errorf("failed to perform request: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return 0, fmt.Errorf("failed to read response: %w", err)
		}

		nonce, err := strconv.ParseUint(strings.TrimSpace(string(body)), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse nonce: %w", err)
		}

		return nonce, nil
	}

	// Validator mode
	if info.NonceReq == nil || info.ValidatorAddr == "" {
		return 0, fmt.Errorf("nonceReq and validatorAddr are required when useIndexer is false")
	}

	conn, err := grpc.NewClient(info.ValidatorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return 0, fmt.Errorf("failed to connect to validator: %w", err)
	}
	defer conn.Close()

	client := helper.NewValidatorNetworkClient(conn)

	response, err := client.Nonce(context.Background(), info.NonceReq)
	if err != nil {
		// If first time
		if strings.Contains(err.Error(), "does not exist") {
			return 1, nil
		}
		return 0, fmt.Errorf("nonce request failed: %w", err)
	}

	return response.GetNonce() + 1, nil // add one from validator to return the nonce the user should currently use
}

func MakeNonceRequest(address string) (*pb.NonceRequest, error) {
	decodedAddr, err := transcode.Base58Decode(address)

	if err != nil {
		return nil, fmt.Errorf("failed to decode address: %w", err)
	}

	return &pb.NonceRequest{WalletAddress: decodedAddr}, nil
}

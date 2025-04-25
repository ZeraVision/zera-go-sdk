package nonce

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/helper"
	"github.com/ZeraVision/zera-go-sdk/transcode"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type NonceInfo struct {
	UseIndexer    bool               // use the ZV indexer if true, false use the validator gRPC api
	NonceReqs     []*pb.NonceRequest // required when useIndexer false
	ValidatorAddr string             // required when useIndexer false
	Addresses     []string           // required when useIndexer true
	IndexerURL    string             // required when useIndexer true
	Authorization string             // required when useIndexer true, Api-Key or Bearer
	Override      []uint64           // optional, if set, this nonce will be used instead of the one from the indexer or validator
}

// GetNonce retrieves the nonce either from the Indexer HTTP API or the Validator gRPC service.
func GetNonce(info NonceInfo, maxRps int) ([]uint64, error) {
	if len(info.Override) > 0 {
		return info.Override, nil
	}

	// Ensure maxRps is valid
	if maxRps <= 1 {
		return []uint64{}, fmt.Errorf("maxRps must be greater than 1")
	}

	// Calculate delay between requests
	delay := time.Second / time.Duration(maxRps-1)

	var (
		nonceRet []uint64
		mu       sync.Mutex
		errChan  = make(chan error, len(info.Addresses)+len(info.NonceReqs)) // Channel to capture errors
	)

	// Create a worker pool channel to limit concurrency
	workerPool := make(chan struct{}, maxRps)

	// Function to process a single address or request
	processNonce := func(addr string, req *pb.NonceRequest, useIndexer bool) {
		defer func() { <-workerPool }() // Release the worker slot when done

		// Wait for rate limiter
		time.Sleep(delay)

		if useIndexer {
			// Indexer mode
			if strings.HasPrefix(addr, "gov_") {
				mu.Lock()
				nonceRet = append(nonceRet, 0)
				mu.Unlock()
				return
			}

			url := fmt.Sprintf("%s/store?requestType=getNextNonce&address=%s", info.IndexerURL, addr)

			req, err := http.NewRequest("POST", url, nil)
			if err != nil {
				errChan <- fmt.Errorf("failed to create request: %w", err)
				return
			}

			req.Header.Add("Target", "indexer")

			// Bearer
			if strings.Contains(info.Authorization, ".") {
				req.Header.Add("Authorization", "Bearer "+info.Authorization)
			} else { // Api Key
				req.Header.Add("Authorization", "Api-Key "+info.Authorization)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				errChan <- fmt.Errorf("failed to perform request: %w", err)
				return
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				errChan <- fmt.Errorf("failed to read response: %w", err)
				return
			}

			nonce, err := strconv.ParseUint(strings.TrimSpace(string(body)), 10, 64)
			if err != nil {
				errChan <- fmt.Errorf("failed to parse nonce: %w", err)
				return
			}

			mu.Lock()
			nonceRet = append(nonceRet, nonce)
			mu.Unlock()
		} else {
			// Validator mode
			if strings.HasPrefix(string(req.WalletAddress), "gov_") {
				mu.Lock()
				nonceRet = append(nonceRet, 0)
				mu.Unlock()
				return
			}

			if info.ValidatorAddr == "" {
				errChan <- fmt.Errorf("validatorAddr is required when useIndexer is false")
				return
			}

			if !strings.Contains(info.ValidatorAddr, ":") {
				info.ValidatorAddr += ":50053"
			}

			conn, err := grpc.Dial(info.ValidatorAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				errChan <- fmt.Errorf("failed to connect to validator: %w", err)
				return
			}
			defer conn.Close()

			client := helper.NewValidatorNetworkApiClient(conn)

			response, err := client.Nonce(context.Background(), req)
			if err != nil {
				// If first time
				if strings.Contains(err.Error(), "does not exist") {
					response = &pb.NonceResponse{Nonce: 0}
				} else {
					errChan <- fmt.Errorf("nonce request failed: %w", err)
					return
				}
			}

			mu.Lock()
			nonceRet = append(nonceRet, response.GetNonce()+1)
			mu.Unlock()
		}
	}

	// Process Indexer mode
	if info.UseIndexer {
		if info.IndexerURL == "" || info.Authorization == "" || len(info.Addresses) < 1 {
			return []uint64{}, fmt.Errorf("indexerURL, authorization, and address are required when useIndexer is true")
		}

		for _, addr := range info.Addresses {
			workerPool <- struct{}{} // Acquire a worker slot
			go processNonce(addr, nil, true)
		}
	}

	// Process Validator mode
	for _, req := range info.NonceReqs {
		workerPool <- struct{}{} // Acquire a worker slot
		go processNonce("", req, false)
	}

	// Wait for all workers to finish
	for i := 0; i < cap(workerPool); i++ {
		workerPool <- struct{}{}
	}
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return nonceRet, nil
}

func MakeNonceRequest(address string) (*pb.NonceRequest, error) {

	if strings.HasPrefix(address, "gov_") {
		return &pb.NonceRequest{WalletAddress: []byte(address)}, nil
	}

	decodedAddr, err := transcode.Base58Decode(address)

	if err != nil {
		return nil, fmt.Errorf("failed to decode address: %w", err)
	}

	return &pb.NonceRequest{WalletAddress: decodedAddr, Encoded: false}, nil
}

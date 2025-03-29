package parts

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strings"
)

type Response struct {
	SupplyInfo struct {
		Parts *big.Int `json:"parts"`
	} `json:"supplyInfo"`
	TokenInfo struct {
		Type string `json:"type"`
	} `json:"tokenInfo"`
}

func GetParts(useIndexer bool, symbol string, endpoint string, authorization string) (*big.Int, error) {
	if symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	if useIndexer {

		if authorization == "" {
			return nil, fmt.Errorf("authorization (api key or bearer token) is required when useIndexer is true")
		}

		// Create the request
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/store?requestType=getContractGlance&symbol=%s", endpoint, symbol), bytes.NewBuffer([]byte{}))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Add headers
		req.Header.Set("Target", "indexer")
		if strings.Contains(authorization, ".") {
			req.Header.Set("Authorization", "Bearer "+authorization)
		} else {
			req.Header.Set("Authorization", "Api-Key "+authorization)
		}

		// Send request
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("API request failed: %v", err)
		}
		defer resp.Body.Close()

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response: %v", err)
		}

		// Parse JSON
		var result Response
		err = json.Unmarshal(body, &result)
		if err != nil {
			if strings.Contains(string(body), "does not exist") {
				return nil, fmt.Errorf("contract with symbol %s does not exist", symbol)
			}
			return nil, fmt.Errorf("failed to parse JSON response: %v", err)
		}

		// Check type
		tokenType := strings.ToLower(result.TokenInfo.Type)
		if tokenType == "sbt" || tokenType == "nft" {
			return nil, fmt.Errorf("%s is %s and has no denomination (always 1 part)", symbol, tokenType)
		}

		// Return parts
		return result.SupplyInfo.Parts, nil
	} else {
		return nil, fmt.Errorf("validator mode is not possible as of network version v.1.1.0")
	}
}

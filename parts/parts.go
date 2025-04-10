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

type PartsInfo struct {
	Symbol     string // contract id
	UseIndexer bool   // true to use indexer, false to use validator
	IndexerUrl string // must be present if UseIndexer is true
	// request pb to val as available
	ValidatorAddr string   // required when UseIndexer false
	Authorization string   // required when useIndexer true, Api-Key or Bearer
	Override      *big.Int // to just specify it
}

func GetParts(partsInfo PartsInfo) (*big.Int, error) {

	if partsInfo.Override != nil {
		return partsInfo.Override, nil
	}

	if partsInfo.Symbol == "" {
		return nil, fmt.Errorf("symbol is required")
	}

	if partsInfo.UseIndexer {

		if partsInfo.Authorization == "" {
			return nil, fmt.Errorf("authorization (api key or bearer token) is required when useIndexer is true")
		}

		// Create the request
		req, err := http.NewRequest("POST", fmt.Sprintf("%s/store?requestType=getContractGlance&symbol=%s", partsInfo.IndexerUrl, partsInfo.Symbol), bytes.NewBuffer([]byte{}))
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %v", err)
		}

		// Add headers
		req.Header.Set("Target", "indexer")
		if strings.Contains(partsInfo.Authorization, ".") {
			req.Header.Set("Authorization", "Bearer "+partsInfo.Authorization)
		} else {
			req.Header.Set("Authorization", "Api-Key "+partsInfo.Authorization)
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
				return nil, fmt.Errorf("contract with symbol %s does not exist", partsInfo.Symbol)
			}
			return nil, fmt.Errorf("failed to parse JSON response: %v", err)
		}

		// Check type
		tokenType := strings.ToLower(result.TokenInfo.Type)
		if tokenType == "sbt" || tokenType == "nft" {
			return nil, fmt.Errorf("%s is %s and has no denomination (always 1 part)", partsInfo.Symbol, tokenType)
		}

		// Return parts
		return result.SupplyInfo.Parts, nil
	} else {
		return nil, fmt.Errorf("validator mode is not possible as of network version v.1.1.0")
	}
}

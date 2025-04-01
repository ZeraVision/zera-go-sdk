package contract

import (
	"fmt"
	"strings"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/helper"
)

// If not specified default false
type RestrictedConfig struct {
	PublicKey      helper.PublicKey // the public key to restrict (if not specified, the contract public key is used). Special r_ keys can be restricted, normal keys can not
	TimeDelay      int64            // can only be put on wallets that don't already exist // TODO verify
	Global         bool             // applies restrictions to wallet accross all contracts //! can only be put on wallets that don't already exist
	UpdateContract bool             // wallets that can upgrade the contract (key weight may create some restrictions regarding other keys)
	Transfer       bool             // is this key allowed to transfer tokens in this token
	Quash          bool             // is this key an authorized key a quash quorum to stop a time delayed transaction
	Mint           bool             // is this key allowed to mint new tokens
	Propose        bool             // is this key allowed to propose new governance proposals within this token
	Vote           bool             // is this key allowed to vote on governance proposals within this token
	Compliance     bool             // is this key allowed to create compliance certificates relating to this token
	ExpenseRatio   bool             // is this key allowed to call in the expense ratio of this token
	CurEquiv       bool             // is this key allowed to publish a self currency equivalent for this token
	Revoke         bool             // is this key allowed to revoke (SBT only) within this token
	KeyWeight      uint32           // lower the weight the higher permissioned it is. Higher permissioned keys can not be removed by lower permissioned keys
}

func CreateRestrictedKeys(config []RestrictedConfig) ([]*pb.RestrictedKey, error) {

	var restricted []*pb.RestrictedKey

	for _, key := range config {

		pubKey, err := helper.GeneratePublicKey(key.PublicKey)
		if err != nil {
			return nil, err
		}

		if !checkRestricted(key.PublicKey) {
			return nil, fmt.Errorf("public key %v is not compatible to be restricted", key.PublicKey)
		}

		// Some of these may not be applicable to your specific contract
		restricted = append(restricted, &pb.RestrictedKey{
			PublicKey:      pubKey,
			TimeDelay:      int64(key.TimeDelay),
			Global:         key.Global,
			UpdateContract: key.UpdateContract,
			Transfer:       key.Transfer,
			Quash:          key.Quash,
			Mint:           key.Mint,
			Propose:        key.Propose,
			Vote:           key.Vote,
			Compliance:     key.Compliance,
			ExpenseRatio:   key.ExpenseRatio,
			CurEquiv:       key.CurEquiv,
			Revoke:         key.Revoke,
			KeyWeight:      key.KeyWeight,
		})
	}

	return restricted, nil
}

func checkRestricted(publicKey helper.PublicKey) bool {

	if publicKey.Single != nil {
		if strings.HasPrefix(*publicKey.Single, "r_") {
			return true
		}
	}

	if publicKey.Multi != nil {
		if publicKey.Multi.HashTokens[0] == helper.RESTRICTED {
			return true
		}
	}

	if publicKey.Inheritence != nil || publicKey.Governance != nil || publicKey.SmartContract != nil {
		return true
	}

	return false // empty -- shouldnt happen
}

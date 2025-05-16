package helper

import (
	"errors"
	"regexp"
	"strconv"

	pb "github.com/ZeraVision/go-zera-network/grpc/protobuf"
	"github.com/ZeraVision/zera-go-sdk/transcode"
)

// Only one of these should be set at a time
type PublicKey struct {
	Single        *string
	Inheritence   *string // ie "$ZRA+0000" inherents all $ZRA+0000 contract keys
	Multi         *MultiKeyHelper
	SmartContract *SmartContractHelper
	Governance    *string // simply specify contractID of token (ie $ZRA+0000)
}

type MultiKeyHelper struct {
	MultiKey   []MultiKey        // each key and its class
	Pattern    [][]MultiPatterns // authentication structure
	HashTokens []HashType        // ie []string{BLAKE, SHA3_256, SHA3_512} = c_a_b
}

type MultiKey struct {
	Class     uint32
	PublicKey string
}

type MultiPatterns struct {
	Class    int32
	Required int32
}

type SmartContractHelper struct {
	Name     string
	Instance uint32
}

func GeneratePublicKey(publicKey PublicKey) (*pb.PublicKey, error) {

	if checkPublicKeyCount(publicKey) != 1 {
		return nil, errors.New("exactly one of the public key options must be set")
	}

	pubKey := &pb.PublicKey{}

	if publicKey.Single != nil {
		_, _, pubKeyByte, _ := transcode.Base58DecodePublicKey(*publicKey.Single)

		pubKey = &pb.PublicKey{
			Single: []byte(pubKeyByte),
		}
	}

	if publicKey.Inheritence != nil {
		// Define the regex pattern for $LETTERS+0000 (4 digits at the end)
		pattern := `^\$[A-Z]+\+\d{4}$`
		matched, err := regexp.MatchString(pattern, *publicKey.Inheritence)
		if err != nil {
			return nil, err
		}

		if !matched {
			return nil, errors.New("inheritence key must be in the format $LETTERS+0000 (with exactly 4 digits at the end)")
		}

		pubKey = &pb.PublicKey{
			Single: []byte(*publicKey.Inheritence),
		}
	}

	if publicKey.Multi != nil {

		pubKey = &pb.PublicKey{
			Multi: &pb.MultiKey{},
		}

		// Public keys in format class_keytype_decodedpublickey
		for _, key := range publicKey.Multi.MultiKey {
			_, _, pubKeyByte, _ := transcode.Base58DecodePublicKey(key.PublicKey)

			pubKey.Multi.PublicKeys = append(pubKey.Multi.PublicKeys, pubKeyByte)
		}

		// Translate embedded array into network format
		for _, pattern := range publicKey.Multi.Pattern { // Outer slice
			classArray := []uint32{}
			requiredArray := []uint32{}

			for _, elements := range pattern { // Inner slice
				classArray = append(classArray, uint32(elements.Class))
				requiredArray = append(requiredArray, uint32(elements.Required))
			}

			// Append the converted pattern to MultiPatterns
			pubKey.Multi.MultiPatterns = append(pubKey.Multi.MultiPatterns, &pb.MultiPatterns{
				Class:    classArray,
				Required: requiredArray,
			})
		}

		// Hash type for multi-key
		for _, hash := range publicKey.Multi.HashTokens {
			pubKey.Multi.HashTokens = append(pubKey.Multi.HashTokens, hash.String())
		}
	}

	if publicKey.SmartContract != nil {
		pubKey = &pb.PublicKey{
			SmartContractAuth: []byte("sc_" + publicKey.SmartContract.Name + "_" + strconv.FormatUint(uint64(publicKey.SmartContract.Instance), 10)),
		}
	}

	if publicKey.Governance != nil {
		pubKey = &pb.PublicKey{
			GovernanceAuth: []byte("gov_" + *publicKey.Governance),
		}
	}

	return pubKey, nil
}

func checkPublicKeyCount(publicKey PublicKey) int {
	count := 0

	if publicKey.Single != nil {
		count++
	}
	if publicKey.Inheritence != nil {
		count++
	}
	if publicKey.Multi != nil {
		count++
	}
	if publicKey.SmartContract != nil {
		count++
	}
	if publicKey.Governance != nil {
		count++
	}

	return count
}

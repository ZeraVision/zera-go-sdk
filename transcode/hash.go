package transcode

import (
	"log"
	"math"

	"golang.org/x/crypto/sha3"

	"github.com/zeebo/blake3"
)

// Blake3 calculates the BLAKE3 hash of a byte slice and returns the hash as a byte slice.
func Blake3(input []byte) []byte {
	hasher := blake3.New()
	_, err := hasher.Write(input)
	if err != nil {
		log.Fatal(err)
	}
	return hasher.Sum(nil)
}

// SHA3256 calculates the SHA3-256 hash of a byte slice and returns the hash as a byte slice.
func SHA3256(input []byte) []byte {
	hasher := sha3.New256()
	_, err := hasher.Write(input)
	if err != nil {
		log.Fatal(err)
	}
	return hasher.Sum(nil)
}

// SHA3512 calculates the SHA3-512 hash of a byte slice and returns the hash as a byte slice.
func SHA3512(input []byte) []byte {
	hasher := sha3.New512()
	_, err := hasher.Write(input)
	if err != nil {
		log.Fatal(err)
	}
	return hasher.Sum(nil)
}

// BLAKE3Variable generates a BLAKE3 hash of the given input with the specified output bit length.
func Blake3Variable(input []byte, outputBitLength int) []byte {
	outputByteLength := int(math.Ceil(float64(outputBitLength) / 8.0))
	hasher := blake3.New()
	_, err := hasher.Write(input)
	if err != nil {
		log.Fatalf("Failed to write to BLAKE3 hasher: %v", err)
	}
	hash := make([]byte, outputByteLength)
	hasher.Digest().Read(hash)

	// Truncate the last few bits if necessary
	if outputBitLength%8 != 0 {
		extraBits := 8 - (outputBitLength % 8)
		hash[outputByteLength-1] &= (0xFF << extraBits)
	}

	return hash
}

// SHAKEVariable generates a SHAKE256 hash of the given input with the specified output bit length.
func ShakeVariable(input []byte, outputBitLength int) []byte {
	outputByteLength := int(math.Ceil(float64(outputBitLength) / 8.0))
	hasher := sha3.NewShake256()
	_, err := hasher.Write(input)
	if err != nil {
		log.Fatalf("Failed to write to SHAKE256 hasher: %v", err)
	}
	hash := make([]byte, outputByteLength)
	_, err = hasher.Read(hash)
	if err != nil {
		log.Fatalf("Failed to read from SHAKE256 hasher: %v", err)
	}

	// Truncate the last few bits if necessary
	if outputBitLength%8 != 0 {
		extraBits := 8 - (outputBitLength % 8)
		hash[outputByteLength-1] &= (0xFF << extraBits)
	}

	return hash
}

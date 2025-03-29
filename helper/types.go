package helper

type KeyType int16

const (
	ED25519 KeyType = 1
	ED448   KeyType = 2
)

type HashType int16

const (
	BLAKE3   HashType = 1
	SHA3_256 HashType = 2
	SHA3_512 HashType = 3
)

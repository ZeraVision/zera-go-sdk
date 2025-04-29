package helper

type KeyType int16

const (
	SPECIAL KeyType = 0 // special key -- ie gov, sc, no sig
	ED25519 KeyType = 1
	ED448   KeyType = 2
	Unknown KeyType = 32767
)

type HashType int16

// There are more -- but these are most common
const (
	RESTRICTED HashType = 0 // not technically a hash, but is used for r_
	BLAKE3     HashType = 1
	SHA3_256   HashType = 2
	SHA3_512   HashType = 3
)

func (h HashType) String() string {
	switch h {
	case BLAKE3:
		return "c"
	case SHA3_256:
		return "a"
	case SHA3_512:
		return "b"
	case RESTRICTED:
		return "r"
	default:
		return "Unknown"
	}
}

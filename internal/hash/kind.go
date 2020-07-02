package hash

const (
	Unknown         = 0x00
	SHA3256         = 0x01
	MurMur3128      = 0x02
	CIVv0SHA2256Raw = 0x03
	CIVv1SHA2256Raw = 0x04
)

func GetKind(hash []byte) (int, error) {
	kind, _, err := parse(hash)
	return kind, err
}

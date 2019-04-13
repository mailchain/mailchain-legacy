package encoding

// DataPrefix used before
func DataPrefix() []byte {
	return []byte{0x6d, 0x61, 0x69, 0x6c, 0x63, 0x68, 0x61, 0x69, 0x6e}
}

const (
	ID        byte = 0x00
	Protobuf  byte = 0x50
	AES256CBC byte = 0x2e // TODO: Not merged yet to multihash
)

package stores

type Transaction struct {
	EnvelopeData []byte
	BlockNumber  int64
	Hash         []byte
}

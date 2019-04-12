package kdf

// StorageKeyDeriver calculates a storage key based on the derivation method
type StorageKeyDeriver interface {
	DeriveKey(options DeriveOpts) ([]byte, error)
}

// DeriveOpts options for deriving a storage key
type DeriveOpts interface {
	KDF() string
}

// DeriveOptionsBuilder used to build DeriveOpts
type DeriveOptionsBuilder func(*DeriveOpts)

package stores

const (
	// SizeSingle is the length of a byte.
	SizeSingle = 1
	// SizeKilobyte is the length of a kilobyte.
	SizeKilobyte = 1 << (10 * 1)
	// SizeMegabyte is the length of a megabyte.
	SizeMegabyte = 1 << (10 * 2)
	// SizeGigabyte is the length of a gigabyte.
	SizeGigabyte = 1 << (10 * 3)
	// SizeTerabyte is the length of a terabyte.
	SizeTerabyte = 1 << (10 * 4)
)

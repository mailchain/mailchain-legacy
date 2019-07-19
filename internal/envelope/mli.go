package envelope

const (
	MLIEmpty     uint64 = 0
	MLIMailchain uint64 = 1
)

// // LocationCode maps the location to the code
// func LocationCode() map[string]uint64 {
// 	return map[string]uint64{
// 		locationMailchain: CodeMailchain,
// 	}
// }

// MLIToAddress maps code to a location
func MLIToAddress() map[uint64]string {
	return map[uint64]string{
		MLIMailchain: mliMailchain,
	}
}

const (
	mliMailchain = "https://mcx.mx"
)

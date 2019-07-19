package msl

// import (
// 	"github.com/mailchain/mailchain/internal/mail"
// 	"github.com/mailchain/mailchain/stores"
// )

// // NewData create data from sent store
// func NewData(sent stores.Sent, hash []byte) (*mail.Data, error) {
// 	store := TypeToCode(sent)
// 	data := &mail.Data{
// 		MsgStoreCode: store,
// 		Hash:         hash,
// 	}
// 	if store == CodeEmpty {
// 		// set schema
// 		// set host
// 	}

// 	return data, nil
// }

// func TypeToCode(sent stores.Sent) uint64 {
// 	switch sent.(type) {
// 	case *stores.SentStore:
// 		return CodeMailchain
// 	default:
// 		return CodeEmpty
// 	}
// }

// const (
// 	SchemaEmpty uint64 = 0x00
// 	SchemaHTTPS uint64 = 0x01
// 	SchemaHTTP  uint64 = 0x02
// )
// const (
// 	CodeEmpty     uint64 = 0x00
// 	CodeMailchain uint64 = 0x01
// )

// // MessageStoreLocation

// // LocationCode maps the location to the code
// func LocationCode() map[string]uint64 {
// 	return map[string]uint64{
// 		locationMailchain: CodeMailchain,
// 	}
// }

// // CodeToLocation maps code to a location
// func CodeToLocation() map[uint64]string {
// 	return map[uint64]string{
// 		CodeMailchain: locationMailchain,
// 	}
// }

// const (
// 	locationMailchain = "https://mcx.mx"
// )

package stores

import "time"

type Header struct {
	// When the message was created, this can be different to the transaction data of the message.
	// example: 12 Mar 19 20:23 UTC
	Date time.Time `json:"date"`
	// The sender of the message
	// example: Charlotte <5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum>
	From string `json:"from"`
	// The recipient of the message
	// To: <4cb0a77b76667dac586c40cc9523ace73b5d772bd503c63ed0ca596eae1658b2@ropsten.ethereum>
	To string `json:"to"`
	// Reply to if the reply address is different to the from address.
	ReplyTo string `json:"reply-to,omitempty"`
	// RekeyTo the address to use when responding.
	RekeyTo string `json:"rekey-to,omitempty"`
	// Unique identifier of the message
	// example: 47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471@mailchain
	MessageID string `json:"message-id"`
	// The content type and the encoding of the message body
	// example: text/plain; charset=\"UTF-8\",
	// 			text/html; charset=\"UTF-8\"
	ContentType string `json:"content-type"`
}

type Message struct {
	// Headers
	Headers Header `json:"headers"`
	// Body of the mail message
	// example: Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac.
	// Sed vestibulum dui lectus, quis eleifend urna mollis eu. Integer dictum metus ut sem rutrum aliquet.
	Body string `json:"body,omitempty"`
	// Subject of the mail message
	// example: Hello world
	Subject    string `json:"subject,omitempty"`
	Status     string `json:"status"`
	StatusCode string `json:"status-code"`
	// Read status of the message
	// example: true
	Read bool `json:"read"`
	// Transaction's block number
	BlockID string `json:"block-id,omitempty"`
	// Transaction's block number encoding type used by the specific protocol
	BlockIDEncoding string `json:"block-id-encoding,omitempty"`
	// Transaction's hash
	TransactionHash string `json:"transaction-hash,omitempty"`
	// Transaction's hash encoding type used by the specific protocol
	TransactionHashEncoding string `json:"transaction-hash-encoding,omitempty"`
}

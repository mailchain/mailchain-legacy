package stores

import "github.com/mailchain/mailchain/internal/pkg/mail"

// Inbox all the actions that support inbox functionality
type Inbox interface {
	DeleteMessageRead(messageID mail.ID) error
	PutMessageRead(messageID mail.ID) error
	GetReadStatus(messageID mail.ID) (bool, error)
}

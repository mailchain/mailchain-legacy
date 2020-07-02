package pinata

import (
	"bytes"
	"net/http"

	"github.com/mailchain/mailchain/encoding"
	"github.com/mailchain/mailchain/internal/mail"
	"github.com/mailchain/mailchain/internal/mli"
	"github.com/pkg/errors"
	"gopkg.in/resty.v1"
)

const pinFileURL = "https://api.pinata.cloud/pinning/pinFileToIPFS"

// NewSent creates a new pinata sen store.
func NewSent(apiKey, apiSecret string) (*Sent, error) {
	return &Sent{
		apiKey:    apiKey,
		apiSecret: apiSecret,
		client:    http.DefaultClient,
	}, nil
}

// Sent handles storing messages in pinata
type Sent struct {
	apiKey    string
	apiSecret string
	client    *http.Client
}

// Key of resource stored.
func (h Sent) Key(messageID mail.ID, contentsHash, msg []byte) string {
	// pref := cid.Prefix{
	// 	Version:  1,
	// 	Codec:    cid.Raw,
	// 	MhType:   mh.SHA2_256,
	// 	MhLength: -1, // default length
	// }
	// pref.Sum(msg)
	return ""
}

// PutMessage stores the message in pinata.
func (h Sent) PutMessage(messageID mail.ID, contentsHash, msg []byte, headers map[string]string) (address, resource string, msgLocInd uint64, err error) {
	// {"IpfsHash":"QmeEYqAaK7UpsffH9fDyPCwB9YH5jf36zDYgeqxhuCvHhh","PinSize":344,"Timestamp":"2020-06-05T22:20:10.110Z"}
	type pinataResponse struct {
		IPFSHash  string `json:"IpfsHash"`
		PinSize   int    `json:"PinSize"`
		Timestamp string `json:"timestamp"`
	}

	type pinataError struct {
		Error string `json:"error"`
	}

	pinRes, err := resty.R().
		SetFileReader("file", encoding.EncodeHexZeroX(messageID), bytes.NewReader(msg)).
		SetFormData(map[string]string{"pinataOptions": `{"cidVersion": 1}`}).
		SetHeader("pinata_api_key", h.apiKey).
		SetHeader("pinata_secret_api_key", h.apiSecret).
		SetResult(&pinataResponse{}).
		SetError(&pinataError{}).
		Post("https://api.pinata.cloud/pinning/pinFileToIPFS")
	if err != nil {
		return "", "", 0, errors.Wrap(err, "pinata: resty error")
	}

	if pinRes.IsError() {
		e := pinRes.Error().(*pinataError)
		return "", "", 0, errors.Errorf("pinata: POST error %s", e.Error)
	}

	out, ok := pinRes.Result().(*pinataResponse)
	if !ok {
		return "", "", 0, errors.Errorf("pinata: POST response invalid")
	}

	return "https://ipfs.io/ipfs/" + out.IPFSHash, out.IPFSHash, mli.IPFS, nil
}

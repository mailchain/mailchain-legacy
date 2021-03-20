package integration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httputil"
	"testing"
	"time"

	"github.com/mailchain/mailchain/crypto"
	"github.com/mailchain/mailchain/encoding"
	"github.com/stretchr/testify/assert"
	"gopkg.in/resty.v1"
)

func apiCheckMessage(t *testing.T, protocol, network, address, checkSubject string) {
	type getResponse struct {
		Messages []message `json:"messages,omitempty"`
	}

	r, err := resty.R().
		SetQueryParams(map[string]string{
			"address":  address,
			"network":  network,
			"protocol": protocol,
		}).
		Get("http://localhost:8080/api/messages")

	if !assert.NoError(t, err) {
		t.FailNow()
	}
	t.Logf("body: %s", r.Body())

	if !assert.Equal(t, 200, r.StatusCode()) {
		t.FailNow()
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/api/response-GET-messages-%s.json", testDir(t), address), r.Body(), 0644); !assert.NoError(t, err) {
		t.FailNow()
	}

	var res getResponse
	if err := json.Unmarshal(r.Body(), &res); !assert.NoError(t, err) {
		t.FailNow()
	}

	if !assert.NotZero(t, res.Messages) {
		t.FailNow()
	}

	if !assert.Equal(t, checkSubject, res.Messages[0].Subject) {
		t.FailNow()
	}
}

func apiFetchMessage(t *testing.T, protocol, network, address string) {
	r, err := resty.R().
		SetQueryParams(map[string]string{
			"address":  address,
			"network":  network,
			"protocol": protocol,
		}).
		Post("http://localhost:8080/api/messages/fetch")

	if !assert.NoError(t, err) {
		t.FailNow()
	}

	if !assert.Equal(t, 200, r.StatusCode()) {
		t.FailNow()
	}
}

func apiSendMessage(t *testing.T, protocol, network string, sendArgs sendArgs, encryptionMethodName, toAddress, fromAddress string, pubKey crypto.PublicKey) string {
	now := time.Now()
	subject := fmt.Sprintf("IT-%s-%d", t.Name(), now.Unix())
	client := resty.New()

	body := map[string]interface{}{
		"content-type":           sendArgs.contentType,
		"encryption-method-name": encryptionMethodName,
		"envelope":               sendArgs.envelope,
		"message": map[string]interface{}{
			"headers": map[string]interface{}{
				"to":   fmt.Sprintf("Charlotte <%s@%s.%s>", toAddress, network, protocol),
				"from": fmt.Sprintf("Sofia <%s@%s.%s>", fromAddress, network, protocol),
			},
			"body":                fmt.Sprintf("Integration test %s. Sending message from Sofia to Charlotte. Time %s", t.Name(), now),
			"subject":             subject,
			"public-key":          encoding.EncodeHexZeroX(pubKey.Bytes()),
			"public-key-encoding": encoding.KindHex0XPrefix,
			"public-key-kind":     pubKey.Kind(),
		},
	}
	req := client.R().SetBody(body)
	res, err := req.Post(fmt.Sprintf("http://localhost:8080/api/messages?protocol=%s&network=%s", protocol, network))
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	dumpReq, err := httputil.DumpRequestOut(req.RawRequest, false)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	reqBody, err := json.Marshal(body)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	dumpReq = append(dumpReq, reqBody...)

	if err := ioutil.WriteFile(fmt.Sprintf("%s/api/request-POST-message-%s.txt", testDir(t), toAddress), dumpReq, 0644); !assert.NoError(t, err) {
		t.FailNow()
	}

	t.Logf("body: %s", res.Body())

	if !assert.Equal(t, 200, res.StatusCode()) {
		t.FailNow()
	}

	return subject
}

func apiGetPublicKey(t *testing.T, address, protocol, network string) *getPublicKeyResponse {
	r, err := resty.R().
		SetResult(getPublicKeyResponse{}).
		SetQueryParams(map[string]string{
			"address":  address,
			"network":  network,
			"protocol": protocol,
		}).
		Get("http://localhost:8080/api/public-key")
	if !assert.NoError(t, err) {
		t.FailNow()
	}
	t.Logf("body: %s", r.Body())

	if !assert.Equal(t, 200, r.StatusCode()) {
		t.FailNow()
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/api/response-GET-public-key-%s.json", testDir(t), address), r.Body(), 0644); !assert.NoError(t, err) {
		t.FailNow()
	}

	return r.Result().(*getPublicKeyResponse)

}

type getPublicKeyResponse struct {
	PublicKey                string   `json:"public-key,omitempty"`
	PublicKeyEncoding        string   `json:"public-key-encoding,omitempty"`
	PublicKeyKind            string   `json:"public-key-kind,omitempty"`
	SupportedEncryptionTypes []string `json:"supported-encryption-types,omitempty"`
}

type message struct {
	BlockID         string `json:"block-id,omitempty"`
	BlockIDEncoding string `json:"block-id-encoding,omitempty"`
	Body            string `json:"body,omitempty"`
	Subject         string `json:"subject,omitempty"`
	Status          string `json:"status,omitempty"`
}

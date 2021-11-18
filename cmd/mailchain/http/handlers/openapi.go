// Copyright 2021 Finobo
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package handlers

//nolint: gofmt
//nolint: lll
//nolint: funlen
func spec() string {
	return `
{
  "consumes": [
    "application/json"
  ],
  "produces": [
    "application/json"
  ],
  "schemes": [
    "https"
  ],
  "swagger": "2.0",
  "info": {
    "description": "All the information needed to talk to the API.\n\nTo raise see anything wrong? Raise an [issue](https://github.com/mailchain/mailchain/issues)",
    "title": "Mailchain API",
    "version": "~mailchain-version~"
  },
  "basePath": "/api",
  "paths": {
    "/addresses": {
      "get": {
        "description": "Get all address that this user has access to. The addresses can be used to send or receive messages.",
        "tags": [
          "Addresses"
        ],
        "summary": "Get addresses.",
        "operationId": "GetAddresses",
        "parameters": [
          {
            "enum": [
              "mainnet",
              "goerli",
              "ropsten",
              "rinkeby",
              "local"
            ],
            "type": "string",
            "example": "goerli",
            "x-go-name": "Network",
            "description": "Network to use when finding addresses.",
            "name": "network",
            "in": "query"
          },
          {
            "enum": [
              "ethereum",
              " substrate",
              " algorand"
            ],
            "type": "string",
            "example": "ethereum",
            "x-go-name": "Protocol",
            "description": "Protocol to use when finding addresses.",
            "name": "protocol",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/GetAddressesResponse"
          },
          "404": {
            "$ref": "#/responses/NotFoundError"
          },
          "422": {
            "$ref": "#/responses/ValidationError"
          }
        }
      }
    },
    "/envelope": {
      "get": {
        "description": "Get envelope\nThis method returns the available envelope types",
        "tags": [
          "Envelope"
        ],
        "summary": "Get Mailchain envelope",
        "operationId": "GetEnvelope",
        "responses": {
          "200": {
            "$ref": "#/responses/GetEnvelopeResponse"
          }
        }
      }
    },
    "/messages": {
      "get": {
        "description": "Check the protocol, network, address combination for Mailchain messages.",
        "tags": [
          "Messages"
        ],
        "summary": "Get Mailchain messages.",
        "operationId": "GetMessages",
        "parameters": [
          {
            "pattern": "0x[a-fA-F0-9]{40}",
            "type": "string",
            "example": "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
            "x-go-name": "Address",
            "description": "Address to use when looking for messages.",
            "name": "address",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              "mainnet",
              "goerli",
              "ropsten",
              "rinkeby",
              "local"
            ],
            "type": "string",
            "example": "goerli",
            "x-go-name": "Network",
            "description": "Network to use when looking for messages.",
            "name": "network",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              "algorand",
              " ethereum",
              " substrate"
            ],
            "type": "string",
            "example": "ethereum",
            "x-go-name": "Protocol",
            "description": "Protocol to use when looking for messages.",
            "name": "protocol",
            "in": "query",
            "required": true
          },
          {
            "type": "boolean",
            "default": false,
            "example": true,
            "x-go-name": "Fetch",
            "description": "Fetch go to the blockchain to retrieve messages",
            "name": "fetch",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int32",
            "default": 15,
            "example": 15,
            "x-go-name": "Limit",
            "description": "Limit the maximum number of message transactions that will be returned. Used for paging.",
            "name": "limit",
            "in": "query"
          },
          {
            "type": "integer",
            "format": "int32",
            "default": 0,
            "example": 0,
            "x-go-name": "Skip",
            "description": "Skip this number of transactions before returning messages. Used for paging.",
            "name": "skip",
            "in": "query"
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/GetMessagesResponse"
          },
          "422": {
            "$ref": "#/responses/ValidationError"
          }
        }
      },
      "post": {
        "description": "Securely send message on the protocol and network specified in the query string to the address.\nOnly the private key holder for the recipient address can decrypt any encrypted contents.\n\nCreate mailchain message\nEncrypt content with public key\nStore message\nEncrypt location\nStore encrypted location on the blockchain.",
        "tags": [
          "Send"
        ],
        "summary": "Send message.",
        "operationId": "SendMessage",
        "parameters": [
          {
            "enum": [
              "mainnet",
              "goerli",
              "ropsten",
              "rinkeby",
              "local"
            ],
            "type": "string",
            "example": "goerli",
            "x-go-name": "Network",
            "description": "Network to use when sending a message.",
            "name": "network",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              "ethereum"
            ],
            "type": "string",
            "example": "ethereum",
            "x-go-name": "Protocol",
            "description": "Protocol to use when sending a message.",
            "name": "protocol",
            "in": "query",
            "required": true
          },
          {
            "description": "Message to send",
            "name": "Body",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/SendMessageRequestBody"
            }
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/StatusOK"
          },
          "404": {
            "$ref": "#/responses/NotFoundError"
          },
          "422": {
            "$ref": "#/responses/ValidationError"
          }
        }
      }
    },
    "/messages/{message_id}/read": {
      "get": {
        "description": "Messages can be either read or unread.",
        "tags": [
          "Messages"
        ],
        "summary": "Message read status.",
        "operationId": "GetRead",
        "responses": {
          "200": {
            "$ref": "#/responses/GetReadResponse"
          },
          "404": {
            "$ref": "#/responses/NotFoundError"
          },
          "422": {
            "$ref": "#/responses/ValidationError"
          }
        }
      },
      "put": {
        "description": "Mark message as read.",
        "tags": [
          "Messages"
        ],
        "summary": "PutRead.",
        "operationId": "PutRead",
        "parameters": [
          {
            "type": "string",
            "x-go-name": "MessageID",
            "description": "Unique id of the message",
            "name": "message_id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/StatusOK"
          },
          "404": {
            "$ref": "#/responses/NotFoundError"
          },
          "422": {
            "$ref": "#/responses/ValidationError"
          }
        }
      },
      "delete": {
        "description": "Mark message as unread",
        "tags": [
          "Messages"
        ],
        "operationId": "DeleteRead",
        "parameters": [
          {
            "type": "string",
            "x-go-name": "MessageID",
            "description": "Unique id of the message",
            "name": "message_id",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/StatusOK"
          },
          "404": {
            "$ref": "#/responses/NotFoundError"
          },
          "422": {
            "$ref": "#/responses/ValidationError"
          }
        }
      }
    },
    "/nameservice/address/{address}/resolve?network={network}\u0026protocol={protocol}": {
      "get": {
        "description": "Get name for supplied address. The name is typically a human-readable value that can be used in place of the address.\nResolve will query the protocol's name service to find the human-readable name for the supplied address.",
        "tags": [
          "NameService"
        ],
        "summary": "Resolve Address Against Name Service",
        "operationId": "GetResolveAddress",
        "parameters": [
          {
            "type": "string",
            "example": "0x4ad2b251246aafc2f3bdf3b690de3bf906622c51",
            "x-go-name": "Address",
            "description": "name to query to get address for",
            "name": "address",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "mainnet",
              "goerli",
              "ropsten",
              "rinkeby",
              "local"
            ],
            "type": "string",
            "example": "goerli",
            "x-go-name": "Network",
            "description": "Network for the name to resolve",
            "name": "network",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "ethereum"
            ],
            "type": "string",
            "example": "ethereum",
            "x-go-name": "Protocol",
            "description": "Protocol for the name to resolve",
            "name": "protocol",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/GetResolveAddressResponse"
          },
          "404": {
            "$ref": "#/responses/NotFoundError"
          },
          "422": {
            "$ref": "#/responses/ValidationError"
          }
        }
      }
    },
    "/nameservice/name/{domain-name}/resolve?network={network}\u0026protocol={protocol}": {
      "get": {
        "description": "Get address for supplied name. The name is typically a human-readable value that can be used in place of the address.\nResolve will query the protocol's name service to find the address for supplied human-readable name.",
        "tags": [
          "NameService"
        ],
        "summary": "Resolve Name Against Name Service",
        "operationId": "GetResolveName",
        "parameters": [
          {
            "type": "string",
            "example": "mailchain.eth",
            "x-go-name": "Name",
            "description": "name to query to get address for",
            "name": "domain-name",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "mainnet",
              "goerli",
              "ropsten",
              "rinkeby",
              "local"
            ],
            "type": "string",
            "example": "goerli",
            "x-go-name": "Network",
            "description": "Network for the name to resolve",
            "name": "network",
            "in": "path",
            "required": true
          },
          {
            "enum": [
              "ethereum"
            ],
            "type": "string",
            "example": "ethereum",
            "x-go-name": "Protocol",
            "description": "Protocol for the name to resolve",
            "name": "protocol",
            "in": "path",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/GetResolveNameResponse"
          },
          "404": {
            "$ref": "#/responses/NotFoundError"
          },
          "422": {
            "$ref": "#/responses/ValidationError"
          }
        }
      }
    },
    "/protocols": {
      "get": {
        "description": "Get all networks for each protocol that is enabled.",
        "tags": [
          "Protocols"
        ],
        "summary": "Get protocols and the networks.",
        "operationId": "GetProtocols",
        "responses": {
          "200": {
            "$ref": "#/responses/GetProtocolsResponse"
          },
          "404": {
            "$ref": "#/responses/NotFoundError"
          },
          "422": {
            "$ref": "#/responses/ValidationError"
          }
        }
      }
    },
    "/public-key": {
      "get": {
        "description": "This method will get the public key to use when encrypting messages and envelopes.\nProtocols and networks have different methods for retrieving or calculating a public key from an address.",
        "tags": [
          "PublicKey"
        ],
        "summary": "Public key from address.",
        "operationId": "GetPublicKey",
        "parameters": [
          {
            "pattern": "0x[a-fA-F0-9]{40}",
            "type": "string",
            "example": "0xde0b295669a9fd93d5f28d9ec85e40f4cb697bae",
            "x-go-name": "Address",
            "description": "Address to to use when performing public key lookup.",
            "name": "address",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              "mainnet",
              "goerli",
              "ropsten",
              "rinkeby",
              "local"
            ],
            "type": "string",
            "example": "goerli",
            "x-go-name": "Network",
            "description": "Network to use when performing public key lookup.",
            "name": "network",
            "in": "query",
            "required": true
          },
          {
            "enum": [
              "ethereum"
            ],
            "type": "string",
            "example": "ethereum",
            "x-go-name": "Protocol",
            "description": "Protocol to use when performing public key lookup.",
            "name": "protocol",
            "in": "query",
            "required": true
          }
        ],
        "responses": {
          "200": {
            "$ref": "#/responses/GetPublicKeyResponse"
          },
          "404": {
            "$ref": "#/responses/NotFoundError"
          },
          "422": {
            "$ref": "#/responses/ValidationError"
          }
        }
      }
    },
    "/version": {
      "get": {
        "description": "Get version of the running mailchain client application and API.\nThis method be used to determine what version of the API and client is being used and what functionality.",
        "tags": [
          "Version"
        ],
        "summary": "Get version",
        "operationId": "GetVersion",
        "responses": {
          "200": {
            "description": "GetVersionResponseBody",
            "schema": {
              "$ref": "#/definitions/GetVersionResponseBody"
            }
          }
        }
      }
    }
  },
  "definitions": {
    "GetAddressesItem": {
      "type": "object",
      "required": [
        "value",
        "encoding",
        "protocol",
        "network"
      ],
      "properties": {
        "encoding": {
          "description": "Encoding method used for encoding the ¬address¬",
          "type": "string",
          "x-go-name": "Encoding",
          "example": "hex/0x-prefix"
        },
        "network": {
          "description": "Network ¬address¬ is available on",
          "type": "string",
          "x-go-name": "Network",
          "example": "mainnet"
        },
        "protocol": {
          "description": "Protocol ¬address¬ is available on",
          "type": "string",
          "x-go-name": "Protocol",
          "example": "ethereum"
        },
        "value": {
          "description": "Address value",
          "type": "string",
          "x-go-name": "Value",
          "example": "0x5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761"
        }
      },
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "GetEnvelopeResponseBodyElement": {
      "description": "GetEnvelopeResponseBodyElement response",
      "type": "object",
      "required": [
        "type",
        "description"
      ],
      "properties": {
        "description": {
          "description": "The envelope description",
          "type": "string",
          "x-go-name": "Description",
          "example": "Private Message Stored with MLI"
        },
        "type": {
          "description": "The envelope type",
          "type": "string",
          "x-go-name": "Type",
          "example": "0x01"
        }
      },
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "GetMessagesResponseHeaders": {
      "type": "object",
      "properties": {
        "content-type": {
          "description": "The content type and the encoding of the message body",
          "type": "string",
          "x-go-name": "ContentType",
          "readOnly": true,
          "example": "text/plain; charset=\\\"UTF-8\\\","
        },
        "date": {
          "description": "When the message was created, this can be different to the transaction data of the message.",
          "type": "string",
          "format": "date-time",
          "x-go-name": "Date",
          "readOnly": true,
          "example": "12 Mar 19 20:23 UTC"
        },
        "from": {
          "description": "The sender of the message",
          "type": "string",
          "x-go-name": "From",
          "readOnly": true,
          "example": "Bob \u003c5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum\u003e"
        },
        "message-id": {
          "description": "Unique identifier of the message",
          "type": "string",
          "x-go-name": "MessageID",
          "readOnly": true,
          "example": "47eca011e32b52c71005ad8a8f75e1b44c92c99fd12e43bccfe571e3c2d13d2e9a826a550f5ff63b247af471@mailchain"
        },
        "rekey-to": {
          "description": "Rekey use this key when responding.",
          "type": "string",
          "x-go-name": "RekeyTo",
          "readOnly": true
        },
        "reply-to": {
          "description": "Reply to if the reply address is different to the from address.",
          "type": "string",
          "x-go-name": "ReplyTo",
          "readOnly": true
        },
        "to": {
          "description": "The recipient of the message",
          "type": "string",
          "x-go-name": "To",
          "readOnly": true
        }
      },
      "x-go-name": "getHeaders",
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "GetMessagesResponseMessage": {
      "type": "object",
      "properties": {
        "block-id": {
          "description": "Transaction's block number",
          "type": "string",
          "x-go-name": "BlockID",
          "readOnly": true
        },
        "block-id-encoding": {
          "description": "Transaction's block number encoding type used by the specific protocol",
          "type": "string",
          "x-go-name": "BlockIDEncoding",
          "readOnly": true
        },
        "body": {
          "description": "Body of the mail message",
          "type": "string",
          "x-go-name": "Body",
          "readOnly": true,
          "example": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante, sit amet ullamcorper dui hendrerit ac."
        },
        "headers": {
          "$ref": "#/definitions/GetMessagesResponseHeaders"
        },
        "read": {
          "description": "Read status of the message",
          "type": "boolean",
          "x-go-name": "Read",
          "readOnly": true,
          "example": true
        },
        "status": {
          "type": "string",
          "x-go-name": "Status",
          "readOnly": true
        },
        "status-code": {
          "type": "string",
          "x-go-name": "StatusCode",
          "readOnly": true
        },
        "subject": {
          "description": "Subject of the mail message",
          "type": "string",
          "x-go-name": "Subject",
          "readOnly": true,
          "example": "Hello world"
        },
        "transaction-hash": {
          "description": "Transaction's hash",
          "type": "string",
          "x-go-name": "TransactionHash",
          "readOnly": true
        },
        "transaction-hash-encoding": {
          "description": "Transaction's hash encoding type used by the specific protocol",
          "type": "string",
          "x-go-name": "TransactionHashEncoding",
          "readOnly": true
        }
      },
      "x-go-name": "getMessage",
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "GetProtocolsProtocol": {
      "type": "object",
      "properties": {
        "name": {
          "description": "in: body",
          "type": "string",
          "x-go-name": "Name"
        },
        "networks": {
          "description": "in: body",
          "type": "array",
          "items": {
            "$ref": "#/definitions/Network"
          },
          "x-go-name": "Networks"
        }
      },
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "GetPublicKeyResponseBody": {
      "description": "GetPublicKeyResponseBody body response",
      "type": "object",
      "required": [
        "public-key",
        "public-key-encoding",
        "public-key-kind",
        "supported-encryption-types"
      ],
      "properties": {
        "public-key": {
          "description": "The public key encoded as per ¬public-key-encoding¬",
          "type": "string",
          "x-go-name": "PublicKey",
          "example": "0x79964e63752465973b6b3c610d8ac773fc7ce04f5d1ba599ba8768fb44cef525176f81d3c7603d5a2e466bc96da7b2443bef01b78059a98f45d5c440ca379463"
        },
        "public-key-encoding": {
          "description": "Encoding method used for encoding the ¬public-key¬",
          "type": "string",
          "x-go-name": "PublicKeyEncoding",
          "example": "hex/0x-prefix"
        },
        "public-key-kind": {
          "description": "Encoding method used for encoding the ¬public-key¬",
          "type": "string",
          "x-go-name": "PublicKeyKind",
          "example": "[\"secp256k1\", \"sr25519\", \"ed25519\"]"
        },
        "supported-encryption-types": {
          "description": "Supported encryption methods for public keys.",
          "type": "array",
          "items": {
            "type": "string"
          },
          "x-go-name": "SupportedEncryptionTypes",
          "example": [
            "aes256cbc",
            "nacl-ecdh",
            "noop"
          ]
        }
      },
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "GetReadResponseBody": {
      "type": "object",
      "required": [
        "read"
      ],
      "properties": {
        "read": {
          "description": "Read",
          "type": "boolean",
          "x-go-name": "Read",
          "example": true
        }
      },
      "x-go-name": "getBody",
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "GetResolveAddressResponseBody": {
      "description": "GetResolveAddressResponseBody body response",
      "type": "object",
      "required": [
        "name"
      ],
      "properties": {
        "name": {
          "description": "The resolved name",
          "type": "string",
          "x-go-name": "Name",
          "example": "mailchain.eth"
        },
        "status": {
          "description": "The RFC1035 status code describing the outcome of the lookup\n\n+ 0 - No Error\n+ 1 - Format Error\n+ 2 - Server Failure\n+ 3 - Non-Existent Domain\n+ 4 - Not Implemented\n+ 5 - Query Refused",
          "type": "integer",
          "format": "int64",
          "x-go-name": "Status",
          "example": 3
        }
      },
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "GetResolveNameResponseBody": {
      "description": "GetResolveNameResponseBody body response",
      "type": "object",
      "required": [
        "address"
      ],
      "properties": {
        "address": {
          "description": "The resolved address",
          "type": "string",
          "x-go-name": "Address",
          "example": "0x4ad2b251246aafc2f3bdf3b690de3bf906622c51"
        },
        "status": {
          "description": "The rFC1035 status code describing the outcome of the lookup\n\n+ 0 - No Error\n+ 1 - Format Error\n+ 2 - Server Failure\n+ 3 - Non-Existent Domain\n+ 4 - Not Implemented\n+ 5 - Query Refused",
          "type": "integer",
          "format": "int64",
          "x-go-name": "Status",
          "example": 3
        }
      },
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "GetVersionResponseBody": {
      "description": "GetVersionResponseBody response",
      "type": "object",
      "required": [
        "version",
        "commit",
        "time"
      ],
      "properties": {
        "commit": {
          "description": "The resolved version commit",
          "type": "string",
          "x-go-name": "VersionCommit"
        },
        "time": {
          "description": "The resolved version release date",
          "type": "string",
          "x-go-name": "VersionDate",
          "example": "2019-09-04T21:59:26Z"
        },
        "version": {
          "description": "The resolved version tag",
          "type": "string",
          "x-go-name": "VersionTag",
          "example": "1.0.0"
        }
      },
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "Network": {
      "type": "object",
      "properties": {
        "id": {
          "type": "string",
          "x-go-name": "ID"
        },
        "name": {
          "type": "string",
          "x-go-name": "Name"
        },
        "nameservice-address-enabled": {
          "type": "boolean",
          "x-go-name": "NameServiceAddressEnabled"
        },
        "nameservice-domain-enabled": {
          "type": "boolean",
          "x-go-name": "NameServiceDomainEnabled"
        }
      },
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "PostMessagesResponseHeaders": {
      "description": "PostHeaders body",
      "type": "object",
      "required": [
        "from",
        "to"
      ],
      "properties": {
        "from": {
          "description": "The sender of the message",
          "type": "string",
          "x-go-name": "From",
          "example": "Bob \u003c5602ea95540bee46d03ba335eed6f49d117eab95c8ab8b71bae2cdd1e564a761@ropsten.ethereum\u003e"
        },
        "reply-to": {
          "description": "Reply to if the reply address is different to the from address.",
          "type": "string",
          "x-go-name": "ReplyTo"
        },
        "to": {
          "description": "The recipient of the message",
          "type": "string",
          "x-go-name": "To"
        }
      },
      "x-go-name": "PostHeaders",
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "PostMessagesResponseMessage": {
      "description": "PostMessage body",
      "type": "object",
      "required": [
        "headers",
        "body",
        "subject",
        "public-key",
        "public-key-encoding"
      ],
      "properties": {
        "body": {
          "description": "Body of the mail message",
          "type": "string",
          "x-go-name": "Body",
          "example": "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Curabitur maximus metus ante,"
        },
        "headers": {
          "$ref": "#/definitions/PostMessagesResponseHeaders"
        },
        "public-key": {
          "description": "Public key of the recipient to encrypt with",
          "type": "string",
          "x-go-name": "PublicKey"
        },
        "public-key-encoding": {
          "description": "Public key Encoding",
          "type": "string",
          "x-go-name": "PublicKeyEncoding"
        },
        "public-key-kind": {
          "description": "Public key kind\nrequire: true",
          "type": "string",
          "x-go-name": "PublicKeyKind"
        },
        "subject": {
          "description": "Subject of the mail message",
          "type": "string",
          "x-go-name": "Subject",
          "example": "Hello world"
        }
      },
      "x-go-name": "PostMessage",
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    },
    "SendMessageRequestBody": {
      "description": "PostRequestBody body",
      "type": "object",
      "required": [
        "message",
        "envelope",
        "encryption-method-name"
      ],
      "properties": {
        "content-type": {
          "description": "Message content-type provided by the client\nrequired: false (default text/plain; charset=\\\"UTF-8\\\")",
          "type": "string",
          "enum": [
            "'text/plain; charset=\\\"UTF-8\\\"'",
            " 'text/html; charset=\\\"UTF-8\\\"'"
          ],
          "x-go-name": "ContentType"
        },
        "encryption-method-name": {
          "description": "Encryption method name",
          "type": "string",
          "enum": [
            "aes256cbc",
            " nacl",
            " noop"
          ],
          "x-go-name": "EncryptionName"
        },
        "envelope": {
          "description": "Envelope that should be used",
          "type": "string",
          "enum": [
            "0x01",
            " 0x50"
          ],
          "x-go-name": "Envelope"
        },
        "message": {
          "$ref": "#/definitions/PostMessagesResponseMessage"
        }
      },
      "x-go-name": "PostRequestBody",
      "x-go-package": "github.com/mailchain/mailchain/cmd/mailchain/http/handlers"
    }
  },
  "responses": {
    "GetAddressesResponse": {
      "description": "GetAddressesResponse Holds the response messages",
      "schema": {
        "type": "array",
        "items": {
          "$ref": "#/definitions/GetAddressesItem"
        }
      }
    },
    "GetEnvelopeResponse": {
      "description": "GetEnvelopeResponse envelope response",
      "schema": {
        "type": "array",
        "items": {
          "$ref": "#/definitions/GetEnvelopeResponseBodyElement"
        }
      }
    },
    "GetMessagesResponse": {
      "description": "GetResponse Holds the response messages",
      "schema": {
        "type": "array",
        "items": {
          "$ref": "#/definitions/GetMessagesResponseMessage"
        }
      }
    },
    "GetProtocolsResponse": {
      "description": "GetProtocolsResponse Holds the response messages",
      "schema": {
        "type": "array",
        "items": {
          "$ref": "#/definitions/GetProtocolsProtocol"
        }
      }
    },
    "GetPublicKeyResponse": {
      "description": "GetPublicKeyResponse public key from address response",
      "schema": {
        "$ref": "#/definitions/GetPublicKeyResponseBody"
      }
    },
    "GetReadResponse": {
      "schema": {
        "$ref": "#/definitions/GetReadResponseBody"
      }
    },
    "GetResolveAddressResponse": {
      "description": "GetResolveAddressResponse address of resolved name",
      "schema": {
        "$ref": "#/definitions/GetResolveAddressResponseBody"
      }
    },
    "GetResolveNameResponse": {
      "description": "GetResolveNameResponse address of resolved name",
      "schema": {
        "$ref": "#/definitions/GetResolveNameResponseBody"
      }
    },
    "GetVersionResponse": {
      "description": "GetVersionResponse version response",
      "schema": {
        "$ref": "#/definitions/GetVersionResponseBody"
      }
    },
    "StatusOK": {
      "description": "StatusOK Description of an StatusOK."
    }
  }
}
`
}

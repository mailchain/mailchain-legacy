# Mailchain

![https://travis-ci.com/mailchain/mailchain](https://img.shields.io/travis/com/mailchain/mailchain.svg?style=for-the-badge)
![https://coveralls.io/github/mailchain/mailchain](https://img.shields.io/coveralls/github/mailchain/mailchain.svg?style=for-the-badge)

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

## Introduction
Mailchain enables blockchain-based email-like messaging with plain or rich text and attachment capabilities. Using blockchain protocols and decentralized storage, Mailchain delivers a simple, secure, messaging experience.

Account holders or owners of a public address often need to communicate in the context of the account or public address. Communication may relate to transactions, account actions or some type of notification.

Examples of use cases for blockchain messaging include invoicing, receipts, contacting users of a service or application developers.

The majority of blockchain protocols provide no standard way to handle messaging. Some account holders have sent messages as an encrypted or unencrypted string included in transaction data. Some applications work around the problem by asking users to link another method of contact (e.g. email address, phone number etc.) to an application. This compromises anonymity by asking users to link or reveal an identity. It also increases exposure to security risks, and relies on additional centralized services.

This proposal outlines how Mailchain gives users the ability to send and receive rich-media HTML messages between public addresses through a simple, email-like interface. All message contents and attachments are encrypted so only the intended recipient of the message can decrypt and view messages.

## Mailchain Message Lifecycle ##

Mailchain is a simple, secure and practical standard which can be implemented across different blockchains. It uses underlying native blockchain protocol capabilities including addressing, immutability, data transmission, and cryptography.

The basic message flow is as follows:

* A user sends a message. The message is encrypted and stored using the recipient’s public key so that only the recipient can locate and decrypt the message contents. The encrypted location of the message is sent by the sender to the recipient in the data field of a transaction.

* To read a message, a recipient uses a corresponding private key to decrypt the location of the message from the data field of a transaction and decrypt the message contents.

---

## Getting started

1. Navigate to https://github.com/mailchain/mailchain/releases/latest to download the latest mailchain build.

### Checksums

To check the application you downloaded matches the build SHA256 checksum.

For example, on a Mac:

``` bash
shasum -a 256 path_to_file/mailchain_...
```

The ouptut should match the checksum found in checksums.txt

## Spec

See https://github.com/mailchain/mailchain-specification for the full requirments that Mailchain is built to address.

### Blockchain Transaction Data

Blockchains store transactions which act as notifications of a message being sent. The location of where the message is stored will be stored in the transaction data ('tx_data'). Each blockchain protocol may have differences in the way data is handled. Mailchain uses a standard format plus any protocol specifics.

*For example: In the case of Ethereum the data will be prefixed with 0x and added to the input field.* 

An explaination of the Ethereum implementation can be count in [internal/pkg/chains/ethereum/README.md](internal/pkg/chains/ethereum/README.md) file.

#### Format

tx_data is formatted as follows `tx_data = encoded(ChainPrefix + Prefix + MultiFormatData)`. All values are stored are bytes that are then `encoded` to be stored on the chain. The `ChainPrefix` allows for any blockchain protocol-specific formatting to be applied.

[MultiFormat](https://github.com/multiformats/multiformats) is used referenced in Mailchain to increase flexibility and conform to existing efforts to build a standard way of expressing format.

`MultiFormatData` is data stored as MultiFormat. It is extensible so future options in how the data is serialized can be accomodated. `MultiFormatData = SerializationType + SerializedData`. Default serialization is [Protocol Buffers](https://developers.google.com/protocol-buffers/) the identifier is `byte{0x50}`. The full spec of the protocol buffer message is described below.

```proto
message Data {
  int32 version = 1;
  bytes encryptedLocation = 2; // MultiFormat
  bytes hash = 3; // MultiFormat
}
```

* `Data.Version` describes the version of the protocol buffer, to support future expansion.

  *Note: protocol buffers support adding fields, the version will need to be incremented in this case. Removing or changing order or type of fields will result in incrementing the version.*

* `Data.encryptionLocation` is stored as MultiFormat to support different encryption methods. `encryptedLocation = EncryptionType + EncryptedData`.
  
  * Default `EncryptionType` is AES-256-CBC
  * `EncryptedData = SentStore + MessageID + EncryptedDataHash`
  
  The `EncryptedDataHash` is used to verify the contents have not been altered.
  *Default: [Murmur3](https://en.wikipedia.org/wiki/MurmurHash)*

* `Data.hash` is stored as MultiFormat to support different hash methods `hash = HashType + HashedData`.

##### Summary of format fields

| Field                          | Value                                                           | Description
| ----------------------------   | --------------------------------------------------------------- | -----------
| ChainPrefix                    | `0x` (Ethereum)          | Optional - dependent on chain requirements
| Prefix                         | `[]byte{0x6d, 0x61, 0x69, 0x6c, 0x63, 0x68, 0x61, 0x69, 0x6e}`  | `6d61696c636861696e` is equal to hexEncoded(`mailchain`).<br/><br/>This field allows the user to identify this transaction is a Mailchain transaction.
| SerializationType              | `byte{0x50}`                                                    | Describes what format the Data is serialized in [MultiFormat](https://github.com/multiformats/multiformats).
| Data.Version                   | `int(1)`                                                        | Maps to a [semantic version](https://semver.org/) of mailchain application.<br/><br/>To allow for extensibility and compatibility. TODO: semantic version mapping
| Data.EncryptedLocation         | `[]byte{MultiFormat....}`                                       | MultiFormat representation of encryption type + encrypted location of the mail message
| Data.EncryptedLocation[0]      | `byte{0x2e}`                                                    | Encrypted type used. *Default: AES256CBC*
| Data.EncryptedLocation[1:]     | `[]byte{EncryptedLocation...}`                                  | Encrypted location of the message
| Data.Hash                      | `[]byte{MultiFormat....}`                                       | MultiFormat representation of encrypted message contents.
| Data.Hash[0]                   | `byte{0x22}`                                                    | Hash function used. *Default: [SHA3-256](https://en.wikipedia.org/wiki/SHA-3)*
| Data.Hash[1:]                  | `[]byte{hash...}`                                               | Hash of plain message contents.

### Storage

Mailchain will support any object storage, including IPFS.

Today, Mailchain supports using S3.

### MultiCipher

Mailchain is intended to support multiple encryption methods (multichipher).

To be supported, the encryption method should use the elliptic curve of the private key + public key pair for the specified chain and be crypto-graphically secure.

#### AES-256-CBC

Advanced Encryption Standard, with a 256 bit key length, operating in the Chain Block Cipher mode (AES-256-CBC).

## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fmailchain%2Fmailchain.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fmailchain%2Fmailchain?ref=badge_large)

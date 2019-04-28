# Mailchain

![Travis](https://img.shields.io/travis/com/mailchain/mailchain.svg?style=for-the-badge)
![Coveralls](https://img.shields.io/coveralls/github/mailchain/mailchain.svg?style=for-the-badge)

Mailchain is an email-like application. Using blockchain encryption and decentralized storage it delivers a simple, secure, email-like experience to blockchain users. Using Mailchain, Ethereum users can now send and receive rich-media HTML messages with attachments via an Ethereum address.

## Getting started

Navigate to xx Download the client

## Privacy and Encryption

Mailchain considers privacy of account holders to be of utmost importance. All message fields will be encrypted to anyone except the recipient. The medium and location of the message is also encrypted.

The Mailchain application will encrypt and decrypt message locations and messages locally on the client. No unintended, un-encrypted information will be transmitted at any time.

Sender & Recipient Privacy Considerations: It is not possible to obfuscate the sender of a message (or transaction sender) without compromising trust (i.e. passing some information to a 3rd party to ask them to send on behalf of another account). As a result, an observer can see that messages have been sent to a recipient and even watch for existence of 2-way conversations.
Risks: Because no other information is revealed, this is considered no more of a risk than what happens when a normal transaction takes place.

## Spec

### Chain Data

Chains store transactions which act as notifications of a message being sent and the location of where the message is stored. Each chain will have differences on the way the data is sent. Mailchain will use the same format plus chain specifics. *Note: In the case of Ethereum the data will be prefixed with 0x and added to the input field.* [MultiFormat](https://github.com/multiformats/multiformats) is a used to throughout Mailchain increase flexibility.

#### Format

Chain data is formatted as follows `ChainData = encoded(ChainPrefix + Prefix + MultiFormatData + ChainSuffix)`. All values are stored are bytes that are then `encoded` to be stored on the chain. The `ChainPrefix` and `ChainSuffix` allow for any chain specific formatting to be applied.

`MultiFormatData` is data stored as MultiFormat and allow future options in how the data is serialized. `MultiFormatData = SerializationType + SerializedData`. Default serialization is [Protocol Buffers](https://developers.google.com/protocol-buffers/) the identifier is `byte{0x50}`. The full spec of the protocol buffer message is described below.

```proto
message Data {
  int32 version = 1;
  bytes encryptedLocation = 2; // MultiFormat
  bytes hash = 3; // MultiFormat
}
```

- `Data.Version` describes the version of the protocol buffer, to support future expansion. *Note: protocol buffers support adding fields, the version will need to be incremented in this case. Removing or changing order or type of fields will result in incrementing the version.*
- `Data.encryptionLocation` is stored as MultiFormat to support different encryption methods. `encryptedLocation = EncryptionType + EncryptedData`.  *Default: encryption method is AES-256-CBC* `EncryptedData = SentStore + MessageID + EncryptedDataHash` The `EncryptedDataHash` is used to verify the contents have not been altered. *Default: [Murmur3](https://en.wikipedia.org/wiki/MurmurHash)* 
- `Data.hash` is stored as MultiFormat to support different hash methods `hash = HashType + HashedData`.

| Field                          | Value                                                           | Description
| ----------------------------   | --------------------------------------------------------------- | -----------
| ChainPrefix                    | `0x` (Ethereum)          | Optional - dependent on chain requirements
| Prefix                         | `[]byte{0x6d, 0x61, 0x69, 0x6c, 0x63, 0x68, 0x61, 0x69, 0x6e}`  | Allows the user to identify this transaction is a Mailchain transaction. Value is `mailchain` hex encoded in bytes
| SerializationType              | `byte{0x50}`                                                    | Describes what format the Data is serialized in [MultiFormat](https://github.com/multiformats/multiformats).
| Data.Version                   | `int(1)`                                                        | Maps to a [semantic version](https://semver.org/) of mailchain application. To allow for extensibility and compatibility. TODO: semantic version mapping
| Data.EncryptedLocation         | `[]byte{MultiFormat....}`                                       | MultiFormat representation of encryption type + encrypted location of the mail message
| Data.EncryptedLocation[0]      | `byte{0x2e}`                                                    | Encrypted type used. *Default: AES256CBC*
| Data.EncryptedLocation[1:]     | `[]byte{EncryptedLocation...}`                                  | Encrypted location of the message
| Data.Hash                      | `[]byte{MultiFormat....}`                                       | MultiFormat representation of encrypted message contents.
| Data.Hash[0]                   | `byte{0x22}`                                                    | Hash function used. *Default: [SHA3-256](https://en.wikipedia.org/wiki/SHA-3)*
| Data.Hash[1:]                  | `[]byte{hash...}`                                               | Hash of plain message contents.
| ChainSuffix                    | `nil` (Ethereum)                                                | Optional - dependent on chain

### Mail message

Messages are encrypted and stored publicly.

#### Storage

Messages are not stored on the blockchain they are stored in a publicly addressable location. The address of the message is stored on the blockchain and encrypted. Only the recipient is able to decrypt the message and determine its location.

#### Encrypted

Different encryption methods are possible

#### Decrypted

### MultiCipher

Mailchain messages and transactions should be encrypted, multiple encryption methods can be supported. The key requirements is that the encryption method support the elliptic curve of the private key + public key pair for the specified chain and be crypto-graphically secure.

#### AES-256-CBC

Advanced Encryption Standard, with a 256 bit key length, operating in the Chain Block Cipher mode (AES-256-CBC).

### Multi hash

### Serialization


## License
[![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fmailchain%2Fmailchain.svg?type=large)](https://app.fossa.io/projects/git%2Bgithub.com%2Fmailchain%2Fmailchain?ref=badge_large)

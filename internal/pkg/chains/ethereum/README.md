# Example Implementation in Ethereum

<!-- @import "[TOC]" {cmd="toc" depthFrom=1 depthTo=6 orderedList=false} -->

Mailchain standardizes messaging in order for Ethereum DApps and Ethereum users to be able to communicate between one another [see Use Cases for examples of why communication is needed].

The message lifecycle for Ethereum is outlined below:

## Message Preparation

1. In the web interface, a user composes a message, including the following fields:

    Field | Description | Example
    | - | - | - |
    To: | The recipient public address | `0xd5ab4ce3605cd590db609b6b5c8901fdb2ef7fe6`
    From: | The sender public address | `0x92d8f10248c6a3953cc3692a894655ad05d61efb`
    Reply-To: | The public address responses should be sent to | `0x92d8f10248c6a3953cc3692a894655ad05d61efb`
    Subject: | The message subject | 
    Body: | The message body |
    Public Key: | The recipient public key (see note) | `0x69d908510e355beb1d5bf2df8129e5b6401e1969891e8016a0b2300739bbb00687055e5924a2fd8dd35f069dc14d8147aa11c1f7e2f271573487e1beeb2be9d0`

    Note: to determine the recipient public key, an existing transaction needs to be sent on the blockchain by that Ethereum account. For example, to determine Bob’s public key, Bob needs to have sent a transaction.

1. Next, a user ‘sends’ the message. The web interface sends the message as a POST request to the Mailchain app.

1. The app adds the following default fields:

    Field | Description | Example
    | - | - | - |
    Date: | The RFC1123 date format | `2019-04-12T18:21:00+01:00`
    Message-id: | A unique message id composed of 64 chars (32 bytes) + `@mailchain` | `002c40fb138807253554afc5161740ca3dade11db7e74e799c9f6091b904277cb9b839393802dc38b8a815615543@mailchain`
    Content-Type: | As per RFC6532 content type for the contents of the message | text/plain; charset="UTF-8" 
    Content-Transfer-Encoding: | The message body is encoded according to this field | `quoted-printable`

1. A hash of the message payload (headers + body etc.) is created (default: `SHA3-256`).

    [returns `message_payload_hash`]

1.  
    1. The message is encrypted using the recipient public key (default: `AES-256-CBC`). This output is a byte array.

    1. The encrypted data byte array is hashed (default murmur3). 

    [returns `encrypted_message_hash`]

1. The encrypted message is uploaded (PUT) to storage with the following attributes:

    File name: `message-id`-`encrypted_message_hash`, e.g. `002c5d4ba47ce66f9e4b1f36f35e50c357aded81dfb9b98a89b8a80d5ca347b2a16f08dc5d37d255378ddcf3380d-220426516c9b`

    [File contents: encrypted bytes]

    The file storage location is returned. E.g. https://mcx.mx/002c5d4ba47ce66f9e4b1f36f35e50c357aded81dfb9b98a89b8a80d5ca347b2a16f08dc5d37d255378ddcf3380d-220426516c9b
    
    [returns `message_location`]

1. The `message_location` is encrypted using the recipient public key (default AES-256-CBC). This is a byte array.

    [returns `encrypted_message_location`]

    1. A protocol buffer (or protobuf) (https://developers.google.com/protocol-buffers) is created containing the following fields:
    
        Field | Type
        - | -
        Version | int32 field for Mailchain versioning
        Encrypted location | `encrypted_message_location`
        Hash | `message_payload_hash`

        [returns `data_protobuf`]

    1. The following fields are then encoded to build the transaction data:
    
        Field | Example
        - | -
        Chain prefix | e.g. `0x` for Ethereum
        Protocol prefix | e.g. `mailchain`
        Multiformatdata | e.g. the `data_protobuf`

    1. The resulting transaction data is then hex-encoded:
        ```
        0x6d61696c636861696e500a82022e808116a34444592018b5b9483...
        ```

        [returns `tx_data`]


## Message Sending - GAS required
Once a message has been encrypted, stored and the resulting `tx_data` prepared, it can be included in a transaction.

1. The Mailchain application creates a transaction with the following fields:

    Field | Details
    -|-
    nonce | The next incremental number in transaction count
    gasPrice | The gas price (at normal rate)
    gasLimit | The estimated required gas
    To | The public address of the recipient
    value | The amount of Eth to include with the message (defaults to 0)
    data | The hex encoded, encrypted location of the Mailchain message (`tx_data` output of Message Preparation).
	
1. The transaction is signed using the sender private key.

1. The transaction is broadcast to the network.

1. The transaction is included in a block.



## Retrieve Messages (GAS NOT required)
The Mailchain application retrieves messages as follows:

1. Retrieve all transactions from the address received transaction history 

1. Identify transactions that contains tx_data beginning `0x6d61696c636861696e...`

1. Decode the message data to get the `data_protobuf`

1. From `data_protobuf`, extract the `encrypted_message_location`

1. Decrypt the `encrypted_message_location` using the recipient private key to obtain the `message_location`

1. Retrieve the encrypted message byte array from the `message_location`

1. Compare the encrypted message byte array with the `encrypted_message_hash` (parsed from the file name in the `message_location`)

1. Decrypt the message using the recipient private key to obtain the `message_payload_hash`

1. Ensure the decrypted contents match the `message_payload_hash` (found in the `data_protobuf`)

1. Return messages. NOTE: messages are returned with a status of ‘ok’ or a description of what went wrong.


## Reading Messages (GAS NOT required)
In the Mailchain web interface inbox, a user can view messages returned by the Mailchain application.

## Key Storage

The Mailchain application handles private key storage using NACL (native go-lang extension). Scrypt
Private keys are used to:

* Decrypt Message locations
* Decrypt messages 
* Sign transactions containing message information

```
comment: // N is the N parameter of Scrypt encryption algorithm, using 256MB
// memory and taking approximately 1s CPU time on a modern processor.
o.N = 1 << 18
// P is the P parameter of Scrypt encryption algorithm, using 256MB
// memory and taking approximately 1s CPU time on a modern processor.
o.P = 1

o.R = 8
o.Len = 32

^^^^ Same as Ethereum ^^^^^

```

## GAS Handling
Gas is spent without authorization. This is considered acceptable risk because values of transactions are hard coded as 0.

The Mailchain application estimates the amount of GAS required to send the transaction and the transaction data at a normal price.

The price for sending a message is usually less than 42000 GAS (2x the cost of a basic transaction).

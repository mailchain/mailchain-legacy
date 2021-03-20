# Integration tests

`private-keys.yaml` contains the private keys used for integration tests. `private-keys.yaml` contains secrets and is not check-in. It needs to be created locally. Below is a sample:


```
algorand:
  testnet:
    charlotte:
      private-key: "words"
      private-key-encoding: "mnemonic/algorand"
      key-type: ed25519
    
    sofia:
      private-key: "words"
      private-key-encoding: "mnemonic/algorand"
      key-type: ed25519

ethereum:
  goerli:
    charlotte:
      private-key: "0x....."
      private-key-encoding: "hex/0x-prefix"
      key-type: ed25519
    
    sofia:
      private-key: "0x....."
      private-key-encoding: "hex/0x-prefix"
      key-type: ed25519      
```

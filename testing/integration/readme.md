# Integration tests

`private-keys.yaml` contains the private keys used for integration tests. `private-keys.yaml` contains secrets and is not check-in. It needs to be created locally. Below is a sample:


```
algorand:
  testnet:
    bob:
      private-key: "words"
      private-key-encoding: "mnemonic/algorand"
      key-type: ed25519
    
    alice:
      private-key: "words"
      private-key-encoding: "mnemonic/algorand"
      key-type: ed25519

ethereum:
  goerli:
    bob:
      private-key: "0x....."
      private-key-encoding: "hex/0x-prefix"
      key-type: ed25519
    
    alice:
      private-key: "0x....."
      private-key-encoding: "hex/0x-prefix"
      key-type: ed25519      
```

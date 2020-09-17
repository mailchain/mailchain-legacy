module github.com/mailchain/mailchain

go 1.13

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/Masterminds/squirrel v1.4.0
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412
	github.com/andreburgaud/crypt2go v0.11.0
	github.com/aws/aws-sdk-go v1.34.6
	github.com/cenkalti/backoff v2.2.1+incompatible
	github.com/centrifuge/go-substrate-rpc-client v2.0.0-alpha.5+incompatible
	github.com/dgraph-io/badger/v2 v2.0.3
	github.com/ethereum/go-ethereum v1.9.19
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/gorilla/mux v1.7.4
	github.com/gtank/merlin v0.1.1
	github.com/gtank/ristretto255 v0.1.2
	github.com/ipfs/go-cid v0.0.6
	github.com/jmoiron/sqlx v1.2.0
	github.com/lib/pq v1.7.0
	github.com/mailchain/go-substrate-rpc-client v2.0.0-alpha.5+incompatible // indirect
	github.com/manifoldco/promptui v0.7.0
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mr-tron/base58 v1.2.0
	github.com/multiformats/go-multihash v0.0.14
	github.com/pierrec/xxHash v0.1.5 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/rubenv/sql-migrate v0.0.0-20200616145509-8d140a17f351
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/afero v1.3.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/syndtr/goleveldb v1.0.1-0.20190923125748-758128399b1d
	github.com/ttacon/chalk v0.0.0-20160626202418-22c06c80ed31
	github.com/urfave/negroni v1.0.0
	github.com/wealdtech/go-ens/v3 v3.4.3
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	gopkg.in/resty.v1 v1.12.0
)

replace github.com/centrifuge/go-substrate-rpc-client v2.0.0-alpha.5+incompatible => github.com/mailchain/go-substrate-rpc-client v2.0.0-alpha.5-mc0.1+incompatible

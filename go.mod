module github.com/mailchain/mailchain

go 1.16

require (
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/Masterminds/squirrel v1.4.0
	github.com/agl/ed25519 v0.0.0-20170116200512-5312a6153412
	github.com/andreburgaud/crypt2go v0.11.0
	github.com/aws/aws-sdk-go v1.34.6
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/centrifuge/go-substrate-rpc-client v2.0.0-alpha.5+incompatible
	github.com/dgraph-io/badger/v2 v2.0.3
	github.com/dgraph-io/ristretto v0.0.2 // indirect
	github.com/ethereum/go-ethereum v1.9.19
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/golang/mock v1.4.4
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/gopherjs/gopherjs v0.0.0-20181103185306-d547d1d9531e // indirect
	github.com/gorilla/mux v1.7.4
	github.com/gtank/merlin v0.1.1
	github.com/gtank/ristretto255 v0.1.2
	github.com/ipfs/go-cid v0.0.6
	github.com/jmoiron/sqlx v1.2.0
	github.com/kr/text v0.2.0 // indirect
	github.com/lib/pq v1.7.0
	github.com/mailchain/go-substrate-rpc-client v2.0.0-alpha.5+incompatible // indirect
	github.com/manifoldco/promptui v0.7.0
	github.com/mattn/go-colorable v0.1.7 // indirect
	github.com/minio/blake2b-simd v0.0.0-20160723061019-3f5f724cb5b1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/mr-tron/base58 v1.2.0
	github.com/multiformats/go-multihash v0.0.14
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pierrec/xxHash v0.1.5 // indirect
	github.com/pkg/errors v0.9.1
	github.com/rs/cors v1.7.0
	github.com/rs/zerolog v1.20.0
	github.com/rubenv/sql-migrate v0.0.0-20200616145509-8d140a17f351
	github.com/sirupsen/logrus v1.6.0
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/afero v1.3.2
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/ttacon/chalk v0.0.0-20160626202418-22c06c80ed31
	github.com/urfave/negroni v1.0.0
	github.com/wealdtech/go-ens/v3 v3.4.3
	golang.org/x/crypto v0.0.0-20200820211705-5c72a883971a
	golang.org/x/net v0.0.0-20200904194848-62affa334b73 // indirect
	golang.org/x/sys v0.0.0-20200905004654-be1d3432aa8f // indirect
	golang.org/x/text v0.3.3 // indirect
	golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1 // indirect
	google.golang.org/protobuf v1.25.0 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/ini.v1 v1.57.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gopkg.in/resty.v1 v1.12.0
	gopkg.in/yaml.v3 v3.0.0-20200615113413-eeeca48fe776 // indirect
)

replace github.com/centrifuge/go-substrate-rpc-client v2.0.0-alpha.5+incompatible => github.com/mailchain/go-substrate-rpc-client v2.0.0-RC6-mc0.1+incompatible

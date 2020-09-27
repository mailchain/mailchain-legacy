# Development Guide

## Running a development version

Mailchain is developed and tested against golang 1.13.x you can confirm this by checking the go version in `.travis.yml`. Use the same go minor version when developing or testing Mailchain from source code.

Fork the repository and clone it to your local environment (see [contributing guidelines](https://github.com/mailchain/mailchain/blob/master/CONTRIBUTING.md#we-use-github-flow)).

To compile and run the package locally, you need to run the following command in place of mailchain command found in the docs:
`go run cmd/mailchain/main.go` + `COMMAND` e.g. `serve`, `account list` etc..

### Examples

#### Add account

To `add account` to your development version

1. Navigate into the directory of the repository
2. Run: `go run cmd/mailchain/main.go account add --protocol=ethereum --key-type="secp256k1" --private-key=YOUR_PRIVATE_KEY`

#### Serve

To `serve` your development version

1. Navigate into the directory of the repository
1. Run: `go run cmd/mailchain/main.go serve`

### To run substrate

For mainnet and beresheet networks:

1. Run the command corresponding to the network you want to send message on `make substrate-mainnet`, `make substrate-beresheet`, or `make substrate-local`. *Note: pull lastest version, if there has been a release it will ask you to continue with new image. Confirm with "y"*
1. [Add keys](#add-account) and [start Mailchain client](#serve).
1. Set protocol to `substrate` and network to the desired option in [Mailchain settings](https://inbox.mailchain.xyz/#/settings).
1. Open [Mailchain inbox](https://inbox.mailchain.xyz/).
1. Send from an SR25519 address to any SR25519 address.

### Inbox vs Inbox Staging

From time to time, there may be pending breaking changes that have not been released. This means that if you use the web inbox, it may not work as desired.

We track this and patch the <https://inbox-staging.mailchain.xyz> to address these teething issues. If you think something may be broken using the regular inbox, try this. If still broken, please [raise an issue](CONTRIBUTING.md#report-bugs-using-githubs-issues).

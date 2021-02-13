# Development Guide

## Running a development version

Development is done using the [Go Programming Language](https://golang.org/).
The version of go is specified in the project's [go.mod](go.mod) file. This document assumes that you have a functioning
environment setup. If you need assistance setting up an environment please visit
the [official Go documentation website](https://golang.org/doc/).

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

Mailchain requires contracts module installed to run. Edgeware has this module installed by default. Mailchain runs against mainnet, beresheet and local networks:

1. Run the command corresponding to the network you want to send message on `make edgeware-mainnet`, `make edgeware-beresheet`, or `make edgeware-local`. *Note: pull lastest version, if there has been a release it will ask you to continue with new image. Confirm with "y"*
2. [Add keys](#add-account) and [start Mailchain client](#serve).
3. Set protocol to `substrate` and network to the desired option in [Mailchain settings](https://inbox.mailchain.xyz/#/settings).
4. Open [Mailchain inbox](https://inbox.mailchain.xyz/).
5. Send from an SR25519 address to any SR25519 address.

### Inbox vs Inbox Staging

From time to time, there may be pending breaking changes that have not been released. This means that if you use the web inbox, it may not work as desired.

We track this and patch the <https://inbox-staging.mailchain.xyz> to address these teething issues. If you think something may be broken using the regular inbox, try this. If still broken, please [raise an issue](CONTRIBUTING.md#report-bugs-using-githubs-issues).

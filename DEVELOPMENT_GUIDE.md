# Development Guide

## Running a development version

Fork the repository and clone it to your local environment (see [contributing guideleines](https://github.com/mailchain/mailchain/blob/master/CONTRIBUTING.md#we-use-github-flow).

To compile and run the package locally, you need to run the following command in place of mailchain command found in the docs:
`go run cmd/mailchain/main.go` + `COMMAND` e.g. `serve`, `account list` etc..

### Examples:
**To `serve` your development version**

1. Navigate into the directory of the repository
1. Run: `go run cmd/mailchain/main.go serve`

**To `add account` using to your development version**

1. Navigate into the directory of the repository
1. Run: `go run cmd/mailchain/main.go account add --protocol=ethereum --key-type="secp256k1" --private-key=YOUR_PRIVATE_KEY`


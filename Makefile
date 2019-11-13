GO := go

all : build
.PHONY: clean test all

clean:
	$(GO) clean
build:
	$(GO) build ./...

test: generate unit-test
unit-test:
	$(GO) test ./...

license: .FORCE
	addlicense -l apache -c Finobo ./internal

proto:
	rm -f ./internal/envelope/*.pb.go
	protoc ./internal/envelope/data.proto -I. --go_out=:.

.PHONY: go-generate
go-generate:
	go generate ./...

.PHONY: generate
generate: go-generate license	

OPENAPIFILE := ./cmd/mailchain/internal/http/handlers/openapi.go
openapi:
	go mod vendor
	rm -rf vendor/github.com/ethereum

	docker run --rm -i \
	-e GOPATH=/go \
	-v $(CURDIR):/go/src/github.com/mailchain/mailchain \
	-w /go/src/github.com/mailchain/mailchain \
	mailchain/goswagger-tool swagger generate spec -b ./cmd/mailchain/internal/http/handlers -o ./docs/openapi/spec.json

	echo "\n""package handlers""\n" > $(OPENAPIFILE)
	echo 'const spec = `' >> $(OPENAPIFILE)
	cat ./docs/openapi/spec.json | sed 's/`/`+"`"+`/g' >> $(OPENAPIFILE)
	echo '`' >>  $(OPENAPIFILE)

	addlicense -l apache -c Finobo $(OPENAPIFILE)
	rm -rf vendor

snapshot:
	docker run --rm --privileged -v $(CURDIR):/go/src/github.com/mailchain/mailchain -v /var/run/docker.sock:/var/run/docker.sock -w /go/src/github.com/mailchain/mailchain mailchain/goreleaser-xcgo goreleaser --snapshot --rm-dist

lint: 
	golangci-lint run --fix

.FORCE:

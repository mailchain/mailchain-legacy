PKGS := github.com/mailchain/mailchain/...
GO := go

all : build
.PHONY: clean test all

clean:
	$(GO) clean
build:
	$(GO) build $(PKGS)

test: generate unit-test
unit-test:
	$(GO) test $(PKGS)

license: .FORCE
	addlicense -l apache -c Finobo ./internal

proto:
	rm -f ./internal/pkg/mail/*.pb.go
	protoc ./internal/pkg/mail/data.proto -I. --go_out=:.

generate:
	sh ./scripts/generate.sh

openapi:
	go mod vendor
	rm -rf vendor/github.com/ethereum
	docker run --rm -i \
	-e GOPATH=/go \
	-v $(CURDIR):/go/src/github.com/mailchain/mailchain \
	-w /go/src/github.com/mailchain/mailchain \
	mailchain/goswagger-tool swagger generate spec -b ./cmd/mailchain/internal/http/handlers -o ./docs/openapi/spec.json

	echo "" >>  ./docs/openapi/spec.json

	echo "package handlers" >  ./cmd/mailchain/internal/http/openapi.go
	echo "" >>  ./cmd/mailchain/internal/http/openapi.go
	echo "// nolint: lll" >>  ./cmd/mailchain/internal/http/openapi.go
	echo 'func spec() string {' >>  ./cmd/mailchain/internal/http/openapi.go
	echo '  return `' >>  ./cmd/mailchain/internal/http/openapi.go
	cat  ./docs/openapi/spec.json >>  ./cmd/mailchain/internal/http/openapi.go
	echo '`' >>  ./cmd/mailchain/internal/http/openapi.go
	echo '}' >>  ./cmd/mailchain/internal/http/openapi.go
	addlicense -l apache -c Finobo ./cmd/mailchain/internal/http/openapi.go	
	rm -rf vendor
	
snapshot:
	docker run --rm --privileged -v $(CURDIR):/go/src/github.com/mailchain/mailchain -v /var/run/docker.sock:/var/run/docker.sock -w /go/src/github.com/mailchain/mailchain mailchain/goreleaser-xcgo goreleaser --snapshot --rm-dist

lint: 
	golangci-lint run --fix

.FORCE:

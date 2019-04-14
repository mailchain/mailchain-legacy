PKGS := github.com/mailchain/mailchain/...
GO := go

all : build
.PHONY: clean test all

clean:
	$(GO) clean
build:
	$(GO) build $(PKGS)

test: genmock unit-test
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
	swagger generate spec -m -b ./internal/pkg/http/rest -o ./docs/openapi/spec.json
	echo "" >>  ./docs/openapi/spec.json

	echo "package rest" >  ./internal/pkg/http/rest/openapi.go
	echo "" >>  ./internal/pkg/http/rest/openapi.go
	echo "// nolint: lll" >>  ./internal/pkg/http/rest/openapi.go
	echo 'func spec() string {' >>  ./internal/pkg/http/rest/openapi.go
	echo '\treturn `' >>  ./internal/pkg/http/rest/openapi.go
	cat  ./docs/openapi/spec.json >>  ./internal/pkg/http/rest/openapi.go
	echo '`' >>  ./internal/pkg/http/rest/openapi.go
	echo '}' >>  ./internal/pkg/http/rest/openapi.go
	addlicense -l apache -c Finobo ./internal/pkg/http/rest/openapi.go

lint: 
	golangci-lint run --fix

.FORCE:

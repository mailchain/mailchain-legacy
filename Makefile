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
	swagger generate spec -m -b ./internal/pkg/http/rest -o ./docs/openapi/spec.json
	echo "" >>  ./docs/openapi/spec.json

	echo "package spec" >  ./internal/pkg/http/rest/spec/openapi.go
	echo "" >>  ./internal/pkg/http/rest/spec/openapi.go
	echo "// nolint: lll" >>  ./internal/pkg/http/rest/spec/openapi.go
	echo 'func spec() string {' >>  ./internal/pkg/http/rest/spec/openapi.go
	echo '\treturn `' >>  ./internal/pkg/http/rest/spec/openapi.go
	cat  ./docs/openapi/spec.json >>  ./internal/pkg/http/rest/spec/openapi.go
	echo '`' >>  ./internal/pkg/http/rest/spec/openapi.go
	echo '}' >>  ./internal/pkg/http/rest/spec/openapi.go
	addlicense -l apache -c Finobo ./internal/pkg/http/rest/spec/openapi.go	

snapshot:
	docker run --rm --privileged -v $PWD:/go/src/github.com/mailchain/mailchain -v /var/run/docker.sock:/var/run/docker.sock -w /go/src/github.com/mailchain/mailchain mailchain/goreleaser-xcgo goreleaser --snapshot --rm-dist

lint: 
	golangci-lint run --fix

.FORCE:

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

genmock:
	sh ./scripts/generate.sh

lint: 
	golangci-lint run --fix

.FORCE:

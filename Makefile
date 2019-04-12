PKGS := github.com/mailchain/mailchain/...
GO := go

all : build
.PHONY: clean test all

clean:
	$(GO) clean
build:
	$(GO) build $(PKGS)

test:
	$(GO) test $(PKGS)

license:
	$(GO) get github.com/google/addlicense
	addlicense -l mit -c Finobo .

proto:
	rm -f ./internal/pkg/mail/*.pb.go
	protoc ./internal/pkg/mail/data.proto -I. --go_out=:.

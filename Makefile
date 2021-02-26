GO := go
VER := latest

all : build
.PHONY: clean test all

help:
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@echo '    help               Show this help screen.'
	@echo '    clean              Remove binaries, artifacts and releases.'
	@echo '    test               Generate Unit test.'
	@echo '    unit-test          Run test.'
	@echo '    build              Build project for current platform.'
	@echo '    go-generation      Open go generate.'
	@echo '    generate           Generate License.'
	@echo '    openapi:           Generate Api'
	@echo ''
	@echo ''

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
	# protoc ./internal/envelope/data.proto -I. --go_out=:.
	docker run --rm -v $(CURDIR):$(CURDIR) -w $(CURDIR) znly/protoc:0.4.0 ./internal/envelope/data.proto -I. --go_out=:.

.PHONY: go-generate clear-docker-images
go-generate:
	go generate ./...

.PHONY: generate
generate: go-generate license	

openapi:
	go mod vendor
	rm -rf vendor/github.com/ethereum
	docker run --rm -i \
	-e GOPATH=/go \
	-v $(CURDIR):/go/src/github.com/mailchain/mailchain \
	-w /go/src/github.com/mailchain/mailchain \
	mailchain/goswagger-tool swagger generate spec -b ./cmd/mailchain/internal/http/handlers -o ./docs/openapi/spec.json

	echo "" >>  ./docs/openapi/spec.json

	echo "package handlers" >  ./cmd/mailchain/internal/http/handlers/openapi.go
	echo "" >>  ./cmd/mailchain/internal/http/handlers/openapi.go
	
	echo "//nolint: gofmt" >> ./cmd/mailchain/internal/http/handlers/openapi.go
	echo "//nolint: lll" >>  ./cmd/mailchain/internal/http/handlers/openapi.go
	echo "//nolint: funlen" >>  ./cmd/mailchain/internal/http/handlers/openapi.go
	echo 'func spec() string {' >>  ./cmd/mailchain/internal/http/handlers/openapi.go
	echo '  return `' >>  ./cmd/mailchain/internal/http/handlers/openapi.go
	cat ./docs/openapi/spec.json | sed 's/`/¬/g' >>  ./cmd/mailchain/internal/http/handlers/openapi.go
	echo '`' >>  ./cmd/mailchain/internal/http/handlers/openapi.go
	echo '}' >>  ./cmd/mailchain/internal/http/handlers/openapi.go	
	gofmt -w -s ./cmd/mailchain/internal/http/handlers/openapi.go
	addlicense -l apache -c Finobo ./cmd/mailchain/internal/http/handlers/openapi.go
	rm -rf vendor

snapshot:
	docker run --rm --privileged -v $(CURDIR):/go/src/github.com/mailchain/mailchain -v /var/run/docker.sock:/var/run/docker.sock -w /go/src/github.com/mailchain/mailchain neilotoole/xcgo goreleaser --snapshot --rm-dist

lint: 
	golangci-lint run --fix

.FORCE:

indexer-database-up:
	go run cmd/indexer/main.go database up --master-postgres-password=mailchain --master-postgres-user=mailchain --indexer-postgres-password=indexer --envelope-postgres-password=envelope --pubkey-postgres-password=pubkey

clear-docker-images:
	docker images -a | grep mailchain | grep latest | awk '{print $$3}'
	- docker rmi -f $$(docker images -a | grep mailchain | grep latest | awk '{print $$3}')

docker-common: clear-docker-images
	docker-compose -f docker-compose.common.yml up --remove-orphans --force-recreate --build

docker-recreate-database: 
	docker-compose -f docker-compose.common.yml down -v

edgeware-mainnet: clear-docker-images
	docker-compose -f docker-compose.common.yml -f docker-compose.edgeware.yml -f docker-compose.edgeware.mainnet.yml pull 
	docker-compose -f docker-compose.common.yml -f docker-compose.edgeware.yml -f docker-compose.edgeware.mainnet.yml up --remove-orphans --force-recreate

edgeware-beresheet: clear-docker-images
	docker-compose -f docker-compose.common.yml -f docker-compose.edgeware.yml -f docker-compose.edgeware.beresheet.yml pull 
	docker-compose -f docker-compose.common.yml -f docker-compose.edgeware.yml -f docker-compose.edgeware.beresheet.yml up --remove-orphans --force-recreate

edgeware-local: clear-docker-images
	docker-compose -f docker-compose.common.yml -f docker-compose.edgeware.yml -f docker-compose.edgeware.local.yml pull 
	docker-compose -f docker-compose.common.yml -f docker-compose.edgeware.yml -f docker-compose.edgeware.local.yml up --remove-orphans --force-recreate
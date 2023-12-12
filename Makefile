MAIN := secretSanta
PKG := `go list -mod=vendor -f {{.Dir}} ./...`

ifeq ($(RACE),1)
	RACEFLAG=-race
endif

build:
	@CGO_ENABLED=0 go build -mod=vendor $(RACEFLAG) -o ${NAME} $(MAIN)

run:
	@echo "Compiling"
	@go run -mod=vendor $(RACEFLAG) $(MAIN) -config=conf.toml -verbose -verbose_sql

mod:
	@go mod tidy
	@go mod vendor
	@git add vendor

test:
	@echo "Running tests"
	@go test -mod=vendor $(RACEFLAG) -coverprofile=coverage.txt -covermode count $(PKG)

lint:
	@golangci-lint run -c .golangci.yml

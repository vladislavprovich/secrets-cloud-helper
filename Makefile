SOURCES ?= $(shell find . -name "*.go" -type f)
BINARY_NAME = go-secretshelper
MAIN_GO_PATH=cmd/cmd.go
COMMIT:=$(shell git describe --tags 2>/dev/null || git rev-parse --short HEAD)
NOW:=$(shell date +"%Y-%m-%d_%H-%M-%S")

all: clean lint build

.PHONY: build
build:
	CGO_ENABLED=0 go build -v -ldflags "-X main.commit=$(COMMIT) -X main.date=$(NOW)" -o dist/${BINARY_NAME} ${MAIN_GO_PATH}

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: lint
lint:
	@for file in ${SOURCES} ;  do \
		golint $$file ; \
	done

.PHONY: test
test:
	@go test -coverprofile=cover.out -coverpkg=./... ./...
	@go tool cover -func=cover.out

.PHONY: gen
gen:
	go generate -v ./...

.PHONY: release
release:
	goreleaser --snapshot --rm-dist

.PHONY: clean
clean:
	rm -rf dist/*
	rm -f cover.out
	go clean -testcache
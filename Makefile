SERVER_PATH := ./cmd/server
CLIENT_PATH := ./cmd/client
BIN_SERVER_PATH := ./bin/server
BIN_CLIENT_PATH := ./bin/client
COVERAGE_FILE := cover.out
NO_MOCKS_COVERAGE_FILE := clean_cover.out
REMOVE_FROM_COVER := internal/mocks|proto|docs
p ?= ./...

BUILD_VERSION := v0.0.1
DATE := $(shell date -u +"%d %b %y %H:%M %z")
COMMIT := $(shell git log --pretty=format:%Creset%s --no-merges -1)
ldflags := -ldflags="\
	-X 'github.com/LekcRg/GophKeeper/internal/buildinfo.BuildVersion=$(BUILD_VERSION)' \
	-X 'github.com/LekcRg/GophKeeper/internal/buildinfo.BuildDate=$(DATE)' \
	-X 'github.com/LekcRg/GophKeeper/internal/buildinfo.BuildCommit=$(COMMIT)'"

all: build-all

run-client:
	go run $(CLIENT_PATH)

run-server:
	go run $(SERVER_PATH) -c=./config/server.yaml

build-server:
	go build -o $(BIN_SERVER_PATH) $(SERVER_PATH)

build-client:
	go build -o $(BIN_CLIENT_PATH) $(CLIENT_PATH)

build:
	make build-server
	make build-client

release-client:
	go build $(ldflags) -o $(BIN_CLIENT_PATH) $(CLIENT_PATH)

release-server:
	go build $(ldflags) -o $(BIN_SERVER_PATH) $(SERVER_PATH)

release:
	make release-client
	make release-server

lint:
	golangci-lint run

cover:
	go test -coverprofile=$(COVERAGE_FILE) ./...
	grep -Ev "$(REMOVE_FROM_COVER)" $(COVERAGE_FILE) > $(NO_MOCKS_COVERAGE_FILE)
	go tool cover -func=$(NO_MOCKS_COVERAGE_FILE)

cover-html:
	go test -coverprofile=$(COVERAGE_FILE) $(p)
	grep -Ev "$(REMOVE_FROM_COVER)" $(COVERAGE_FILE) > $(NO_MOCKS_COVERAGE_FILE)
	go tool cover -html=$(NO_MOCKS_COVERAGE_FILE)

make swag:
	swag init -g ./cmd/server/main.go

betteralign:
	betteralign -apply -test_files ./...

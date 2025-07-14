SERVER_PATH := ./cmd/server
CLIENT_PATH := ./cmd/client
COVERAGE_FILE := cover.out
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

build-all:
	go build -o $(CLIENT_PATH)/client $(CLIENT_PATH)
	go build -o $(SERVER_PATH)/server $(SERVER_PATH)

build-server:
	go build -o $(SERVER_PATH)/server $(SERVER_PATH)

release-client:
	go build $(ldflags) -o $(CLIENT_PATH)/client $(CLIENT_PATH)

release-server:
	go build $(ldflags) -o $(SERVER_PATH)/server $(SERVER_PATH)

release:
	make release-client
	make release-server

lint:
	golangci-lint run

cover:
	go test -coverprofile=$(COVERAGE_FILE) ./...
	go tool cover -func=$(COVERAGE_FILE)

cover-html:
	go test -coverprofile=$(COVERAGE_FILE) $(p)
	go tool cover -html=$(COVERAGE_FILE)

make swag:
	swag init -g ./cmd/server/main.go

betteralign:
	betteralign -apply -test_files ./...

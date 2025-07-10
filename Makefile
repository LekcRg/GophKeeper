SERVER_PATH := ./cmd/server
CLIENT_PATH := ./cmd/client

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

lint:
	golangci-lint run

betteralign:
	betteralign -apply -test_files ./...

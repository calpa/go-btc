BIN_NAME=go-btc
BIN_DIR=bin

.PHONY: all run build test tidy fmt

all: build

run:
	go run .

build:
	mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(BIN_NAME) .

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	go fmt ./...

.PHONY: build test lint fmt vet clean all

BUILD_DIR := bin
BINARY := c2b

build:
	go build -o $(BUILD_DIR)/$(BINARY) ./cmd/c2b

test:
	go test ./... -v -count=1

lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .

vet:
	go vet ./...

clean:
	rm -rf $(BUILD_DIR)

all: fmt vet lint test build

.PHONY: build

all: run

run: build
	./build/blockchain

build:
	GOOS=linux go build -o build/blockchain ./cmd/blobkchain

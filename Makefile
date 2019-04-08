.PHONY: build

build:
	go build -a -o ./build/fetch ./cmd/fetch/

test:
	go test ./... -race -v

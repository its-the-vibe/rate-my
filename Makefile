.PHONY: build test lint run

build:
	go build -o rate-my .

test:
	go test ./...

lint:
	go vet ./...
	gofmt -l .

run:
	go run main.go

.PHONY: test build run fmt
fmt:
	gofmt -w ./cmd ./internal

test:
	go test ./...

build:
	go build -o bin/rooommetrics ./cmd/rooommetrics

run:
	go run ./cmd/rooommetrics

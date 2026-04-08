APP_NAME := shipdeck

.PHONY: build run test lint dev clean

build:
	go build -o bin/$(APP_NAME) ./cmd/shipdeck

run:
	air

test:
	go test ./...

lint:
	gofmt -w ./cmd ./internal

dev:
	air

clean:
	rm -rf bin

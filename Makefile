.PHONY: all build run test test-race test-cover lint

all: build run

build:
	docker-compose build

run:
	docker-compose up

test:
	go test -v -timeout=1m ./...

test-race:
	go test -race -v -timeout=1m ./...

test-cover:
	rm -f all.coverage.out
	go test -race -v -timeout=1m \
        -coverprofile=all.coverage.out \
        -coverpkg=./... $$(go list ./...|grep -v cmd)

lint:
	docker pull golangci/golangci-lint:latest
	docker run -v`pwd`:/workspace -w /workspace \
        golangci/golangci-lint:latest golangci-lint run ./...

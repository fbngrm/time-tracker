.PHONY: all build test test-race test-cover lint

all: build

build:
	mkdir -p bin
	go build -o bin/time-tracker \
		-ldflags "-X main.version=$${VERSION:-$$(git describe --tags --always --dirty)}" \
        ./cmd/time-tracker/main.go

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

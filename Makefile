# Tell make to treat all targets as phony
.PHONY: all build test clean run vet install build-unix

# Go parameters
BINARY_NAME=cf-plus-mac
BINARY_UNIX=cf-plus
GO_VERSION=1.10

all: install build

build:
	@go build -o $(BINARY_NAME) -v

# Cross compilation
build-unix: install
	docker run -e GOPATH=/app/ \
		-e CGO_ENABLED=0 \
		--rm  \
		-v "$(PWD):/target" \
		-v "$(PWD):/app/src/github.com/ukayani/cloudformation-plus" \
		-w /target  \
		golang:1.10 \
		go build \
		-o cf-plus

test:
	@go test -v ./...

install:
	dep ensure
clean:
	@go clean
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

vet:
	@go vet
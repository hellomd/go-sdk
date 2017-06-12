GOFILES = $(shell find . -name '*.go' -not -path './vendor/*')
GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

default: build

build:
	go build -o dist/server

run: build
	./dist/server

test: test-all

test-short:
	ENV=test go test -short -v $(GOPACKAGES)

test-all:
	ENV=test go test -v $(GOPACKAGES)

FROM golang:latest

WORKDIR /go/src/github.com/hellomd/go-sdk
ADD . .
ENTRYPOINT go test -v $(go list ./... | grep -v /vendor/)
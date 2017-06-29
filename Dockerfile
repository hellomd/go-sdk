FROM golang:latest
WORKDIR /go/src/github.com/hellomd/middlewares
ADD . .
ENTRYPOINT go test -v $(go list ./... | grep -v /vendor/)
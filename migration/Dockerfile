FROM golang:1.21.0-alpine

RUN apk --no-cache --update --upgrade add ca-certificates gcc libc-dev

ENV GO111MODULE on
ENV CGO_ENABLED 1
ENV GOOS linux
ENV GOARCH amd64

WORKDIR /go/migrations
ADD . /go/migrations
RUN go build -o migration .
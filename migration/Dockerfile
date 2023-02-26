FROM golang:1.19.5-alpine

RUN apk --no-cache --update --upgrade add ca-certificates

ENV GO111MODULE on
ENV CGO_ENABLED 1
ENV GOOS linux
ENV GOARCH amd64

WORKDIR /go/migrations
ADD . /go/migrations
RUN go build -o migration .
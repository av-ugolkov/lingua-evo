FROM golang:1.22.0-alpine

RUN apk --no-cache --update --upgrade add go

ENV GO111MODULE on
ENV CGO_ENABLED 1
ENV GOOS linux
ENV GOARCH amd64

WORKDIR /migration
ADD . /migration
RUN go build -o migration .
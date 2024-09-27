FROM golang:1.22.0-alpine as builder

RUN apk --no-cache --update --upgrade add go

ENV GO111MODULE on
ENV CGO_ENABLED 1
ENV GOOS linux
ENV GOARCH amd64

WORKDIR /migration
ADD . /migration
RUN go build -o migration .

FROM alpine:3.19
LABEL key="Lingua Migration"

RUN --mount=type=cache,target=/var/cache/apk apk --update --upgrade add make bash

COPY --from=builder ./migration ./

ENTRYPOINT ./migration
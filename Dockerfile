FROM golang:1.19.5-alpine as builder

RUN apk --no-cache --update --upgrade add git make

WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go make build

FROM alpine:3.16
MAINTAINER Lingua-evo

RUN apk --no-cache --update --upgrade add curl

WORKDIR .
COPY configs/dev.yaml configs/dev.yaml
COPY --from=0 . .
CMD ["/app/main"]
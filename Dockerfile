FROM golang:1.19.5-alpine as builder

RUN apk --no-cache --update --upgrade add git make

WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go make build

FROM alpine:3.16
MAINTAINER Lingua-evo

ARG config_dir
ARG config_file

RUN apk --no-cache --update --upgrade add curl

WORKDIR .
COPY ./configs/${config_dir}/${config_file}.yaml /configs/${config_dir}/${config_file}.yaml
COPY --from=0 . .
CMD ["/app/main"]
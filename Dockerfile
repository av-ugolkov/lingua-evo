FROM golang:1.19.5-alpine as builder

RUN apk --no-cache --update --upgrade add git make

WORKDIR /app
COPY . .
RUN --mount=type=cache,target=/go make build

FROM alpine:3.16
MAINTAINER Lingua-evo

ARG config_dir

RUN apk --no-cache --update --upgrade add curl

WORKDIR .
COPY ./configs/${config_dir}/server_config.yaml /configs/server_config.yaml
COPY --from=0 . .

EXPOSE 5000

CMD ["/app/main"]
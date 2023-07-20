FROM golang:1.20.3-alpine as builder

RUN apk --no-cache --update --upgrade add git make

WORKDIR /build
COPY ./main .

FROM alpine:3.16
LABEL key="Lingua-evo"

ARG config_dir

RUN apk --no-cache --update --upgrade add curl

WORKDIR /
COPY ./configs/${config_dir}/server_config.yaml /configs/server_config.yaml
COPY ./web /web
COPY --from=0 . .

EXPOSE 5000

CMD ["/build/main"]
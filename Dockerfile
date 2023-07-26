FROM golang:1.20.3-alpine as builder

RUN apk --no-cache --update --upgrade add git make

WORKDIR /build
COPY . .
RUN --mount=type=cache,target=/go make build

FROM alpine:3.16
LABEL key="Lingua-evo"

ARG config_dir

RUN apk --no-cache --update --upgrade add curl

WORKDIR /lingua-evo
COPY /configs/${config_dir}/server_config.yaml ./configs/server_config.yaml
COPY /static/web ./static/web
COPY --from=builder ./build/main ./app/

EXPOSE 5000

WORKDIR /lingua-evo/app/

CMD ["/lingua-evo/app/main"]
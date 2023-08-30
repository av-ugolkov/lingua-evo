FROM golang:1.21.0-alpine as builder

RUN apk --no-cache --update --upgrade add git make

WORKDIR /build
COPY . .
RUN --mount=type=cache,target=/go make build

FROM alpine:3.18
LABEL key="Lingua-evo"

ARG config_dir

RUN apk --no-cache --update --upgrade add curl

WORKDIR /lingua-evo
COPY /configs/${config_dir}/server_config.yaml ./configs/server_config.yaml
COPY /static/web ./static/web
COPY --from=builder ./build/main ./cmd/

EXPOSE 5000

WORKDIR /lingua-evo/cmd/

ENTRYPOINT ["./main"]
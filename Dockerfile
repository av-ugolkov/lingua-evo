FROM golang:1.22.0-alpine as builder

RUN apk --no-cache --update --upgrade add git make ca-certificates openssl

WORKDIR /build
COPY . .
RUN --mount=type=cache,target=/go make build

FROM alpine:3.19
LABEL key="Lingua-evo"

ARG config_dir

RUN apk --update --upgrade add ca-certificates git

WORKDIR /lingua-evo
COPY /configs/${config_dir}/server_config.yaml ./configs/server_config.yaml
COPY --from=builder ./build/cmd/main ./

EXPOSE 5000

WORKDIR /lingua-evo/
CMD ["./main"]
FROM golang:1.22.0-alpine as builder

RUN --mount=type=cache,target=/var/cache/apk apk --no-cache --update --upgrade add git make ca-certificates openssl

WORKDIR /build
COPY . .
RUN --mount=type=cache,target=/go make build

FROM alpine:3.19
LABEL key="Lingua Evo"

ARG config_dir
ARG cert_path
ARG epsw

RUN --mount=type=cache,target=/var/cache/apk apk --update --upgrade add ca-certificates git bash

WORKDIR /lingua-evo

COPY /configs/${config_dir}.yaml ./configs/server_config.yaml
COPY --from=builder ./build/cmd/main ./
COPY --from=root /${cert_path}/certificate.crt ./cert/
COPY --from=root /${cert_path}/private.key ./cert/

EXPOSE 5000

ENV env_epsw=${epsw}

WORKDIR /lingua-evo/
ENTRYPOINT ./main -epsw=$(echo ${env_epsw})
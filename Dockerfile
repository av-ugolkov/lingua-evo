FROM golang:1.22.0-alpine as builder

RUN --mount=type=cache,target=/var/cache/apk apk --no-cache --update --upgrade add git make ca-certificates openssl

WORKDIR /build
COPY . .
RUN --mount=type=cache,target=/go make build

FROM alpine:3.19
LABEL key="Lingua Evo"

ARG config_dir
ARG public_cert
ARG private_cert
ARG epsw
ARG branch
ARG commit

LABEL git.branch=$branch
LABEL git.commit=$commit

RUN --mount=type=cache,target=/var/cache/apk apk --update --upgrade add ca-certificates git bash

WORKDIR /lingua-evo

COPY /configs/${config_dir}.yaml ./configs/server_config.yaml
COPY --from=builder ./build/cmd/main ./
COPY --from=root ${public_cert} ./cert/
COPY --from=root ${private_cert} ./cert/

EXPOSE 5000

ENV env_epsw=${epsw}

WORKDIR /lingua-evo/
ENTRYPOINT ./main -epsw=$(echo ${env_epsw})
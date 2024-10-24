FROM golang:1.23.2-alpine as builder

RUN --mount=type=cache,target=/var/cache/apk apk --no-cache --update --upgrade add git make ca-certificates openssl

WORKDIR /build
COPY . .
RUN --mount=type=cache,target=/go make build

FROM alpine:3.20.3
LABEL key="Lingua Evo"

ARG config_dir
ARG public_cert
ARG private_cert
ARG epsw
ARG jwts
ARG pg_psw
ARG redis_psw
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
ENV env_jwts=${jwts}
ENV env_pg_psw=${pg_psw}
ENV env_redis_psw=${redis_psw}

WORKDIR /lingua-evo/
ENTRYPOINT ./main -epsw=$(echo ${env_epsw}) -jwts=$(echo ${env_jwts}) -pg_psw=$(echo ${env_pg_psw}) -redis_psw=$(echo ${env_redis_psw})
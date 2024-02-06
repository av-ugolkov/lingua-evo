FROM golang:1.21.0-alpine as builder

RUN apk --no-cache --update --upgrade add git make

WORKDIR /build
COPY . .
RUN --mount=type=cache,target=/go make build

FROM scratch
LABEL key="Lingua-evo"

ARG config_dir

WORKDIR /lingua-evo
COPY /configs/${config_dir}/server_config.yaml ./configs/server_config.yaml
COPY --from=builder ./build/cmd/main ./

EXPOSE 5000

WORKDIR /lingua-evo/

ENTRYPOINT ["./main"]
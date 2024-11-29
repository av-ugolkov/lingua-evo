#!/bin/bash

run(){
    cfg="./configs/server_config.yaml"
    epsw="$(cat email.env | grep EPSW | cut -d "=" -f2)"
    jwts="$(cat jwt.env | grep JWTS | cut -d "=" -f2)"
    pg_psw="$(cat db.env | grep PG_PSW | cut -d "=" -f2)"
    redis_psw="$(cat db.env | grep REDIS_PSW | cut -d "=" -f2)"

    go run ./cmd/main.go -config=${cfg} -epsw=${epsw} -jwts=${jwts} -pg_psw=${pg_psw} -redis_psw=${redis_psw}
}

"$@"
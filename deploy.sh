#!/bin/bash

dev() {
  BRANCH="$(cut -d "/" -f2 <<< "$(git rev-parse --abbrev-ref HEAD)")"
  COMMIT="$(git $dir rev-parse HEAD)"

  epsw="$(cat .env | grep EPSW | cut -d "=" -f2)"
  jwts="$(cat .env | grep JWTS | cut -d "=" -f2)"
  pg_psw="$(cat .env | grep PG_PSW | cut -d "=" -f2)"
  redis_psw="$(cat .env | grep REDIS_PSW | cut -d "=" -f2)"

  BRANCH=${BRANCH} \
  COMMIT=${COMMIT} \
  EPSW=${epsw} \
  JWTS=${jwts} \
  PG_PSW=${pg_psw} \
  REDIS_PSW=${redis_psw} \
  docker compose -p lingua-evo-dev -f deploy/docker-compose.dev.yml up --build --force-recreate    
}

database() {
  BRANCH="$(git rev-parse --abbrev-ref HEAD)"
  
  pg_psw="$(cat .env | grep PG_PSW | cut -d "=" -f2)"

  PG_PSW=${pg_psw} \
  BRANCH=${BRANCH} \
  docker compose -p lingua-evo-dev -f deploy/docker-compose.dev.yml up redis postgres migration --build --force-recreat    
}

database_down() {
  BRANCH="$(git rev-parse --abbrev-ref HEAD)"
  
  pg_psw="$(cat .env | grep PG_PSW | cut -d "=" -f2)"

  PG_PSW=${pg_psw} \
  BRANCH=${BRANCH} \
  CMD=down \
  docker compose -p lingua-evo-dev -f deploy/docker-compose.dev.yml up redis postgres migration --build --force-recreate    
}

release() {
  clear

  if test -d /home/lingua-logs; then
    echo "Directory [ lingua-logs ] exists."
  else
    mkdir /home/lingua-logs
  fi

  if test -d /home/lingua-dumps; then
    echo "Directory [ lingua-dumps ] exists."
  else
    mkdir /home/lingua-dumps
  fi

  echo "Do you want to create a backup for database? [y/n]"
  read ans
  if [ "$ans" = "y" ]; then
    ./backup.sh
  fi

  BRANCH="$(cut -d "/" -f2 <<< "$(git rev-parse --abbrev-ref HEAD)")"
  echo "Do you want to deploy from [ $BRANCH ]? [y/n]"
  read ans
  if [ "$ans" = "n" ]; then
    echo "Type branch name which you want yo use"
    read BRANCH
  fi

  CURRENT_BRANCH="$(cut -d "/" -f2 <<< "$(git rev-parse --abbrev-ref HEAD)")"
  if [ "$BRANCH" != "$CURRENT_BRANCH" ]; then
    echo "$(git checkout $BRANCH)"
  fi

  echo "$(git fetch)"
  echo "$(git pull)"

  COMMIT="$(git $dir rev-parse HEAD)"

  epsw="$(cat .env | grep EPSW | cut -d "=" -f2)"
  jwts="$(cat .env | grep JWTS | cut -d "=" -f2)"
  pg_psw="$(cat .env | grep PG_PSW | cut -d "=" -f2)"
  redis_psw="$(cat .env | grep REDIS_PSW | cut -d "=" -f2)"

  EPSW=${epsw} \
  JWTS=${jwts} \
  PG_PSW=${pg_psw} \
  REDIS_PSW=${redis_psw} \
  BRANCH=${BRANCH} \
  COMMIT=${COMMIT} \
  docker compose -p lingua-evo -f deploy/docker-compose.yml up --build --force-recreate
}

"$@"
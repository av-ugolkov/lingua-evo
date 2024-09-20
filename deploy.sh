#!/bin/bash

dev() {
    BRANCH="$(git rev-parse --abbrev-ref HEAD)"
    COMMIT="$(git $dir rev-parse HEAD)"
    echo $(
        BRANCH=${BRANCH} \
        COMMIT=${COMMIT} \
        docker compose -p lingua-evo-dev -f deploy/docker-compose.dev.yml up --build --force-recreate
    )    
}

database() {
    BRANCH="$(git rev-parse --abbrev-ref HEAD)"
    echo $(
        BRANCH=${BRANCH} \
        docker compose -p lingua-evo-dev -f deploy/docker-compose.dev.yml up redis postgres migration --build --force-recreate
    )    
}

database_down() {
    BRANCH="$(git rev-parse --abbrev-ref HEAD)"
    echo $(
        BRANCH=${BRANCH} \
        CMD=down \
        docker compose -p lingua-evo-dev -f deploy/docker-compose.dev.yml up redis postgres migration --build --force-recreate
    )    
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

    BRANCH="$(git rev-parse --abbrev-ref HEAD)"
    echo "Do you want to deploy [ $BRANCH ]? [y/n]"
    read ans
    if [ "$ans" = "n" ]; then
        echo "Type branch name which you want yo use"
        read BRANCH
    fi
    CURRENT_BRANCH="$(git rev-parse --abbrev-ref HEAD)"
    if [ "$BRANCH" != "$CURRENT_BRANCH" ]; then
        echo "$(git checkout $BRANCH)"
    fi

    EPSW="$(cat .env | xargs)"
    
    echo "$(git fetch)"
    echo "$(git pull)"

    COMMIT="$(git $dir rev-parse HEAD)"
    
    echo $(
        EPSW=${EPSW} \
        BRANCH=${BRANCH} \
        COMMIT=${COMMIT} \
        docker compose -p lingua-evo -f deploy/docker-compose.yml up --build --force-recreate
    )
}

"$@"
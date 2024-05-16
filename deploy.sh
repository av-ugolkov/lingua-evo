#!/bin/bash

clear
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
echo "$(make run.docker $EPSW)"

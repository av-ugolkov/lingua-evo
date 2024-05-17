#!/bin/bash

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
echo "$(make run.docker $EPSW)"

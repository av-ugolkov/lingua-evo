#!/bin/bash

clear

docker exec -t -i postgres-base /bin/bash -c "
    cd /home/dump
    pg_dump -U lingua pg-lingua-evo > backup_"$(date +'%d_%m_%Y_%H_%M_%S')".dump"

#!/bin/bash

docker exec -ti postgres-base /bin/bash -c "
    if test -d /home/dumps; then
        echo \"Directory [ dumps ] exists.\"
    else
        mkdir /home/dumps
    fi
    cd /home/dumps

    pg_dump -U lingua pg-lingua-evo > backup_\"$(date +'%d_%m_%Y_%H_%M_%S')\".dump
"

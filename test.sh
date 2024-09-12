#!/bin/bash

# handlers
ADDR=localhost:7070
REGISTER_URL=$ADDR/register
READY_URL=$ADDR/player/ready
START_GAME_URL=$ADDR/admin/start_game


result=$(curl --location $REGISTER_URL --header 'Content-Type: text/plain' --data '{
    "username": "admin",
    "password": "admin"
}' -s)

ADMIN_TOKEN=$(echo $result | jq -r ".token")

result=$(curl --location $REGISTER_URL --header 'Content-Type: text/plain' --data '{
    "username": "hey1",
    "password": "hey"
}' -s)

PLAYER1_TOKEN=$(echo $result | jq -r ".token")

result=$(curl --location $REGISTER_URL --header 'Content-Type: text/plain' --data '{
    "username": "hey2",
    "password": "hey"
}' -s)

PLAYER2_TOKEN=$(echo $result | jq -r ".token")

curl --location $START_GAME_URL --header "Authorization: Bearer $ADMIN_TOKEN" -s &

function kill_game_pids {
    start_game_curl_pids=$(ps aux | grep curl | grep $START_GAME_URL | grep -v grep | awk '{print $2}')
    for pid in $start_game_curl_pids; do kill $pid; done
}
trap kill_game_pids EXIT

result=$(curl --location $READY_URL --header "Authorization: Bearer $PLAYER1_TOKEN" -s)
result=$(curl --location $READY_URL --header "Authorization: Bearer $PLAYER2_TOKEN" -s)

sleep 10

#!/bin/ash

# link upload dir to nginx for indexing
ln -s /uploads /nginx/data/bin

# start redis-server
redis-server /redis.conf &
sleep 1

# start binder
binder &

# start nginx
nginx

#!/bin/sh
export REDIS_HOST=`ping -c 1 codis-proxy.marathon.mesos | head -1 | awk -F "(" '{print $2}' | awk -F ")" '{print $1}'`
nginx -g "daemon off; error_log /dev/stderr info;"

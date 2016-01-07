#!/bin/bash

sed -i s#NGINX_SERVER_NAME#$NGINX_SERVER_NAME# /etc/nginx/conf.d/backends.conf
sed -i s#BACKENDS_LINK#$BACKENDS_LINK# /etc/nginx/conf.d/backends.conf

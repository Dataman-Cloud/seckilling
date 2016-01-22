# warehouse

The app warehouse is the backends system of seckilling in Dataman. The warehouse consists of mysql, redis(codis), django, nginx. In the depolyement, the mysql and redis containers are required in advanced, then you need build the images of django and nginx and start the containers.

## deployment by hands

Start the mysql container,
```
# docker run -d --name db -e MYSQL_ROOT_PASSWORD=111111 -e MYSQL_DATABASE=warehouse -p 3306:3306 mysql:5.6
# mysql -uroot -p111111 -h127.0.0.1 warehouse
```

Start the redis container,
```
# docker run -d --name cache -p 6379:6379 redis:latest
# redis-cli -h 127.0.0.1 -p 6379
```

Build the warehouse app and run it
```
# cd /opt/seckilling/warehouse/backends
# docker build -t warehouse_app:latest .
# cd ..
# docker run -d --name app -p 8000:8000 --link=db --link=cache -e DB_HOST=db -e REDIS_HOST=cache -e REDIS_PORT=6379 -v /opt/warehouse/backends:/code warehouse_app:latest
```

Build the nginx web and run it
```
# cd /opt/seckilling/warehouse/nginx
# docker build -t warehouse_nginx:latest .
# cd ..
# docker run -d --name nginx -p 80:80 --link=app -e BACKENDS_LINK="http://app:8000" -e NGINX_SERVER_NAME=127.0.0.1 warehouse_nginx:latest
```

## Setup by Makefile

Clone the repo
```
# git clone git@github.com:Dataman-Cloud/seckilling.git
# cd seckilling/warehouse/

```

Start the containers by docker-compose
```
# make init
```

Migrate the database warehouse
```
# make init-db
```

Create the superuser
```
# make create-superuser
```

## Clean up the ENV
```
# make cleanup
```

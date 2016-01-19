##秒杀哦服务部署流程

###部署go服务
  镜像地址: testregistry.dataman.io/seckill/gatev:0.1
  启动命令: 
  ```
    docker run \
    -v gate-conf.sample.yaml:/etc/seckilling/gate-conf.yaml \
    p 8090:8090 \
    --name seckilling-gate \
    testregistry.dataman.io/seckill/gatev:0.1 
  ```
  注: 可以不挂在配置文件使用默认配置, 也可以等container启动后进入container内部修改配置文件然后重启容器
  配置说明:
  ```
  ---
  logLevel: "DEBUG"   // 日志级别
  port: ":8090"       // 服务端口
  cache:
    host: "123.59.58.58" // redis 地址
    port: 5506     // redis 端口
    password: ""  // redis 密码 (一般不填)
    db: 0    // db (一般不填)
    poolsize: 100  // redis 连接池容量
  ```

###部署nginx服务
  镜像地址: testregistry.dataman.io/nginx-proxy-1.9.7.1
  启动命令: 
  ```
    PWD=$(pwd)
    NG_PREFIX=/opt/openresty/nginx
    MODE=prod
    docker run \
    -v $PWD/conf/nginx.conf:$NG_PREFIX/conf/nginx.conf \
    -p 9200:80 \
    -e REDIS_HOST="123.59.61.172" \
    -e REDIS_PORT=19000 \
    -e TOKEN_COOKIE="DM_SK_UID" \
    -e SALT_OFFSET=180000 \
    -e COUNTER_BATCH=5 \
    -d \
    --name resty \
    testregistry.dataman.io/nginx-proxy-1.9.7.1:1.0 && docker logs -f resty
  ```
  参数说明: REDIS_HOST REDIS_PORT 为redis的地址和端口
            COUNTER_BATCH 为nginx二级计数的步长
            TOKEN_COOKIE 为cookie的key值
            SALT_OFFSET为salt失效时间
            除了redis的两个参数其他都可以使用默认值
  nginx参数说明:
  ```
    upstream sk.server{
        least_conn;

        server 192.168.1.104:8090; # 改成go服务的IP
        keepalive 10;
    }
    
    location = /api/v1/seckill {
        limit_except GET {
            deny all;
        }
  
  	default_type application/json;
        access_by_lua_file lua/access_seckill.lua;
        header_filter_by_lua_file lua/filter_seckill.lua;
        proxy_pass http://192.168.1.104:8090/api/v1/seckill;   # IP用go程序的IP替换
    }
  ```

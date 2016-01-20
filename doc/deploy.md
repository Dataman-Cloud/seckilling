##秒杀服务部署流程

###部署go服务
  镜像地址: testregistry.dataman.io/seckill/gate:v0.1
  启动命令: 
  ```
    docker run \
    -p 8090:8090 \
    -e INIT_MODEL=1 \
    -e HOST="localhost" \
    -e PORT=":8090" \
    -e LOG_LEVEL="DEBUG" \
    -e CACHE_HOST="123.59.61.172" \
    -e CACHE_PORT=19000 \
    -e CACHE_POOLSIZE=100 \
    -d \
    --name seckilling-gate \
    --net host \
    testregistry.dataman.io/seckill/gate:v0.1
  ```
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
  镜像地址: testregistry.dataman.io/nginx-proxy-1.9.7.1:1.0
  启动命令: 
  ```
    PWD=$(pwd)
    NG_PREFIX=/opt/openresty/nginx
    MODE=prod
    docker run \
    -e REDIS_HOST="123.59.61.172" \
    -e REDIS_PORT=19000 \
    -e TOKEN_COOKIE="DM_SK_UID" \
    -e SALT_OFFSET=180000 \
    -e COUNTER_BATCH=5 \
    -d \
    --name resty \
    -net host \
    testregistry.dataman.io/nginx-proxy-1.9.7.1:1.0 && docker logs -f resty
  ```
  参数说明: REDIS_HOST REDIS_PORT 为redis的地址和端口
            COUNTER_BATCH 为nginx二级计数的步长
            TOKEN_COOKIE 为cookie的key值
            SALT_OFFSET为salt失效时间
            除了redis的两个参数其他都可以使用默认值

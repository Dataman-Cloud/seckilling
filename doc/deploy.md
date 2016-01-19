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

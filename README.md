# seckilling
seckilling platform

# request count limit
采用两级基于计数器的请求数限制, 第一级在Redis中，第二级在Nginx Worker，计数使用后增机制，每个worker默认分配100个请求数，处理完之后增加到第一级计数器，这种机制最多超额100*workers请求到后端，但不会由于Nginx崩溃而少处理请求，超额的请求量由后端拦截, 兼顾性能和请求数控制。

# API
异常状态，http status != 200 ：
```json
{
  "code": 1
}
```

- `GET /api/v1/setevent?id=xxx&effectOn=xxx&duration=xxx&batch=xxx` 管理平台控制活动开始开关.可以主动调用这个接口来开启一个新的活动,
  设置开一一个新的活动需要设置活动的id,开始时间,持续时间以及nginx计数的步长(可以不设使用默认值).管理平台调用这个接口的同时需要向redis中插入相同的数据方便nginx重启时加载数据.再部署多个nginx时管理平台需要轮询分别向每一个nginx发送创建活动的通知.如果设置成功则返回当前缓存的所有活动,否则返回异常状态 <br/>
  ```json
  {
    "time": 1452071210514,
    "events": [
      {
        "desc": "",
        "effectOn": 1451893345000,
        "id": 1,
        "duration": 600000
      },
      {
        "desc": "",
        "effectOn": 1451888328000,
        "id": 2,
        "duration": 240000
      },
      {
        "desc": "",
        "effectOn": 1451888568000,
        "id": 3,
        "duration": 240000
      }
    ]
  }
  ```


- `GET /api/v1/events`  H5加载后第一次请求，通过不同的营销活动ID获取即将进行的活动列表(Cached in Nginx)，其中包含ID、开始时间、描述和资源图片；非200返回显示`无活动`页面<br/>
```json
{
  "time": 1452071210514,
  "events": [
    {
      "desc": "",
      "effectOn": 1451893345000,
      "id": 1,
      "duration": 600000
    },
    {
      "desc": "",
      "effectOn": 1451888328000,
      "id": 2,
      "duration": 240000
    },
    {
      "desc": "",
      "effectOn": 1451888568000,
      "id": 3,
      "duration": 240000
    }
  ]
}
```
- `GET /api/v1/event?id=xxx` 活动页刷新，返回生效时间和服务器时间，客户端进行倒计时，活动开始前xxxm会附加本次秒杀按钮的唯一salt，如果客户端没有salt，需要在进入倒计时xxxm区间是自动刷新获取salt；非200返回显示`活动已结束`页面<br/>
```json
{
  "effectOn": 1451893345000,
  "time": 1452071319261,
  "salt": "69221910-C1C2-4BAD-9F79-F06C8D231209"
}
```

- `GET /api/v1/seckill?id=xxx&phone=xxx&salt=xxx` 秒杀api，进行计数，先来先得，多于限额的用户判断为失败，命中的请求服务器会同时设置cookie DM_SK_UID作为token, 并记录到redis；非200返回显示`秒杀失败`页面<br/>
```json
{
  "coupon": "69221910-C1C2",
}
```

- `GET /api/v1/coupon?phone=xxx` 发送优惠码到手机；非200返回显示`活动已结束`页面<br/>
```json
{
  "coupon": "69221910-C1C2",
}
```

# Redis Key Formats

* `events`: Events list
* `event:<eid>`: Event info hash
    - `id`: Event ID
    - `effectOn`: UTC timestamp second
    - `duration`: Lifetime in seconds for this events
    - `desc`: Description for this event
    - ...
* `sn:<eid>`: Serial Numbers in this event, Sorted Sets
    - `score`: DB index for this Serial Number
    - `element`: Serial Number
* `cur_eid`: Current Event ID
* `count:<eid>`: Delivered Count for this event
* `tr:<eid>:<sn>`: Order result hash for `sn` in event `eid`
    - `phone`: user cell phone number
    - ...

# design
![Design] (doc/design.jpg)

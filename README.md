# seckilling
seckilling platform

# request count limit
采用两级基于计数器的请求数限制, 第一级在Redis中，第二级在Nginx Worker，计数使用后增机制，每个worker默认分配100个请求数，处理完之后增加到第一级计数器，这种机制最多超额100*workers请求到后端，但不会由于Nginx崩溃而少处理请求，超额的请求量由后端拦截, 兼顾性能和请求数控制。

# API
/v1/api

1. GET /:cid/events  H5加载后第一次请求，通过不同的营销活动ID获取即将进行的活动列表(Cached in Nginx)，其中包含ID、开始时间、描述和资源图片<br/>
[{id, effectOn, time, description, resources}]

2. GET /events/:id 活动页刷新，返回生效时间和服务器时间，客户端进行倒计时，活动开始前10m会附加本次秒杀按钮的唯一URL，如果客户端没有URL，需要在进入倒计时10m区间是自动刷新获取URL<br/>
{effectOn, time, url, phone}

3. PUT /events/:id/URL 秒杀api，进行计数，先来先得，多于限额的用户判断为失败，命中的请求服务器会同时设置cookie DM_SK_UID作为token并记录到redis<br/>
{status} 

4. PUT /events/:id/phone 手机号码提交<br/>
{status}

5. PUT /events/:id/order 短信验证码提交并生成订单，记录到redis<br/>
{status, coupons}

# design
![Design] (doc/design.jpg)

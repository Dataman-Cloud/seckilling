### seckill 设计实现说明文档

 * 事件流程 request -> proxy -> queue -> order -> push -> ...
 
### proxy
  proxy 主要用来限流, 添加response计数机制, 通过response计数第一次筛选秒杀请求, 只有在response计数到达规定的最大数目之前的请求才会被发送到queue, 超过最大请求数限制的请求将会被重定向到CDN上的秒杀结束界面.
  
### queue
  queue服务主要提供秒杀页面操作的接口和将秒杀请求的流量转到kafka和redis
  设计实现流程:
* 创建一个新的秒杀活动: </br>
    (1). 设置本次活动开始时间, 活动名称, proxy限流数量, 商品数目等信息作为一条新的记录插入到table event中, 生成活动ID </br>
    (2). 将本次秒杀活动的商品信息录入到table prod_inst中在插入时需要加上该商品的活动ID和序列号加上 </br>
    (3). 在redis中加入一个计数器key值为活动ID, 最大值为本次活动提供的商品数目 </br>
    (4). 前端页面开始倒计时
* 购买(门票)接口
  接收到proxy转发过来的请求之后将用户的唯一标示符插入到kafka和redis, redis以用户为唯一标识存入一个hash  *{"status":0}* 处理完成之后跳转到等待支付页面.
  
### order
  order服务主要对kafka中的ticket进行核验以及进行库存更新进一步确认秒杀成功的用户</br>
  (1) 从redis中获取商品目前剩余数目, 如果库存为0, 核验MySQL中的库存信息如果确认商品已经卖完(或需要终止该活动), 提示queue活动已经结束</br>
  (2) 从kafka中取出一张ticket </br>
  (3) 查询该门票在redis中的状态(如果不存在或者status不为0则当做无效门票处理) </br>
  (4) 生成一条order记录包含ticket, 活动ID, 商品剩余库存等信息插入到table order中 </br>
  (5) 将redis中ticket的状态置为1 </br>
  (6) redis中商品计数器-1 </br>
* 注: 如果活动过程中redis异常或活动异常需要回复, 以MySQL中table order中的库存数量为准重新开始计数  

### push
    push服务对接电商的支付接口, 向电商push订单
    





 

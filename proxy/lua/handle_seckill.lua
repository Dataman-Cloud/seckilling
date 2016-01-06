local config = require "config"
local ck = require "resty.cookie"
local cookie, err = ck:new()
if not cookie then
    ngx.log(ngx.CRIT, "can't create cookie ", err)
    ngx.exit(ngx.HTTP_NOT_ACCEPTABLE)
    return
end
local token = cookie:get(config.tokenCookie)

local redisc = require "redisc"
local redis = redisc:new()
local res, err = redis:hget("tk:"..token, "status")

local uuid = require "uuid4"
local coupon = uuid.getUUID()
ngx.say(string.format('{"coupon":"%s"}', coupon))

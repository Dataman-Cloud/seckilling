local args = ngx.req.get_uri_args()
local config = require "config"
local ck = require "resty.cookie"
local cookie, err = ck:new()
if not cookie then
    ngx.log(ngx.CRIT, "can't create cookie ", err)
    ngx.exit(ngx.HTTP_NOT_ACCEPTABLE)
    return
end

local redisc = require "redisc"
local redis = redisc:new()
local token = args[config.tokenCookie]
local cjson = require "cjson"
print("get args ", cjson.encode(args), " get token ", token)
local res, err = redis:hget("tk:"..token..args.id , "status")

local uuid = require "uuid4"
local coupon = uuid.getUUID()
ngx.say(string.format('{"coupon":"%s"}', coupon))


local uri = ngx.re.sub(ngx.var.uri, "^/v1/api/events/(.*)", "$1", "o")
local cjson = require "cjson";
--redis
local redis = require "redisc"
local red = redis:new()

-- cookie
local ck = require "resty.cookie"
local cookie, err = ck:new()
if not cookie then
    ngx.log(ngx.ERR, err)
    return
end
local field, err = cookie:get("DM_SK_UID")
if not field then
    -- uuid
    local uuid = require("uuid")
    local uuidstr = uuid.generate_random()
    local ok, err = cookie:set({
        key = "DM_SK_UID", value = uuidstr
    })
    ok, err = red:set("DM_SK_UID", uuidstr)
    if not ok then
        ngx.say("can't save token to redis: ", err)
    end
end

local time, terr = red:hget("event:" .. uri, "time")
local effectOn, terr = red:hget("event:" .. uri, "effectOn")
tab1 = {}
tab1["time"] = time
tab1["effectOn"] = effectOn
tab1["serverTime"] = os.date("%Y-%m-%d %H:%M:%S", os.time())
ngx.say(cjson.encode(tab1))

local uri = ngx.re.sub(ngx.var.uri, "^/v1/api/events/(.*)", "$1", "o")
--redis
local redis = require "redis"
local red = redis:new()
red:set_timeout(1000) -- 1 sec
local ok, err = red:connect(addr, port)
if not ok then
    ngx.say("failed to connect: ", err)
    return
end

-- cookie
local ck = require "cookie"
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
--tab1["serverTime"] = os.time()
ngx.say(cjson.encode(tab1))
--[[if res == nil then
    ngx.say("not found event ID: ", uri, err, res)
    return
end]]

--[[local pattern = "(%d+)-(%d+)-(%d+)T(%d+):(%d+):(%d+)"
local runyear, runmonth, runday, runhour, runminute, runseconds = res:match(pattern)
local convertedTimestamp = os.time({year = runyear, month = runmonth, day = runday, hour = runhour - 8, min = runminute, sec = runseconds})
local ts = os.time() + 1000 * 60
if ts > convertedTimestamp then
    ngx.say("alreay start")
end
ngx.say(convertedTimestamp)]]
-- ngx.say(ngx.re.sub(ngx.var.uri, "^/v1/api/events/(.*)",, os.time() "$1", "o"))

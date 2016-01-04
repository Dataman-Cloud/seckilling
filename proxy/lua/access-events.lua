--local uri = ngx.re.sub(ngx.var.uri, "^/v1/api/(.*)/events", "$1", "o")
local cjson = require "cjson";
local uri = "events"

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
    local uuid = require("resty.uuid")
    local uuidstr = uuid.generate_random()
    local ok, err = cookie:set({
        key = "DM_SK_UID", value = uuidstr
    })
    ok, err = red:set("DM_SK_UID", uuidstr)
    if not ok then
        ngx.say("can't save token to redis: ", err)
    end
end


-- get events
local civities, err = red:lrange(uri, 0, -1)
if not civities then
    ngx.say("can't found events list ", err)
    return
end
tab1 = {}
local x = 0
for i, res in ipairs(civities) do
    local id, ierr = red:hget("event:" .. res, "id")
    local effectOn, eerr = red:hget("event:" .. res, "effectOn")
    local time, terr = red:hget("event:" .. res, "time")
    local description, derr = red:hget("event:" .. res, "description")
    local resources, rerr = red:hget("event:" .. res, "resources")
    local tab2 = {}
    tab2["id"] = id
    tab2["effectOn"] = effectOn
    tab2["time"] = time
    tab2["description"] = description
    tab2["resources"] = resources
    tab1[i] = tab2
end
ngx.say(cjson.encode(tab1))

local cache = ngx.shared.scache
local json, err = cache:get("events")
if not json then
    ngx.log(ngx.ERR, "events not found", err)
    ngx.exit(ngx.HTTP_NOT_FOUND)
else 
    ngx.say(string.format([[
    {
        "time":%d,
        "events": %s
    }
    ]], ngx.now() * 1000, json))
end

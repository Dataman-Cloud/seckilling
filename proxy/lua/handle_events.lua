local cache = ngx.shared.scache
local val, err = cache:get("events")
if not val then
    ngx.log(ngx.ERR, "events not found", err)
    ngx.exit(ngx.HTTP_NOT_FOUND)
else 
    ngx.say(val)
end

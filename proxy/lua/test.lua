
local cache = ngx.shared.scache
local events = cache:get("events")
print(events)

ngx.say(events)

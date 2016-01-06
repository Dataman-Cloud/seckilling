local config = require "config"
local args = ngx.req.get_uri_args()
local id = args.id

local cache = ngx.shared.scache
local effectOn, err = cache:get('eeo:'..id)
local duration, err = cache:get('ed:'..id)
print("========eeo ", effectOn, " === ed ", duration)
if not (effectOn and duration) then
    ngx.exit(ngx.HTTP_NOT_FOUND)
else 
    local now = ngx.now() * 1000
    local offset = effectOn - now
    if offset <= config.saltOffset and now < effectOn + duration then
        local cache = ngx.shared.scache
        local salt, err = cache:get('salt:'..id)
        if not salt then
            local uuid = require "uuid4"
            salt = uuid.getUUID()
            cache:set("salt:"..id, salt)
        end
        ngx.say(string.format('{"e":%d, "t":%d, "s": "%s"}', effectOn, now, salt))
    else
        ngx.say(string.format('{"e":%d, "t":%d}', effectOn, now))
    end
end

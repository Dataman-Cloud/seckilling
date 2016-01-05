local config = require "config"
local args = ngx.req.get_uri_args()
local id = args.id

local cache = ngx.shared.scache
local val, err = cache:get('eeo:'..id)
val = 1451995870448 --val * 1000
if not val then
    ngx.exit(ngx.HTTP_NOT_FOUND)
else 
    local now = ngx.now() * 1000
    if val - now <= config.saltOffset then
        local cache = ngx.shared.scache
        local salt, err = cache:get('salt')
        if not salt then
            local uuid = require "uuid4"
            salt = uuid.getUUID()
            cache:set("salt", salt)
        end
        ngx.say(string.format('{"e":%d, "t":%d, "s": "%s"}', val, now, salt))
    else
        ngx.say(string.format('{"e":%d, "t":%d}', val, now))
    end
end

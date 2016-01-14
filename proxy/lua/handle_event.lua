local config = require "config"
local args = ngx.req.get_uri_args()
local id = args.id
local constant = require "constant"

local cache = ngx.shared.scache
local effectOn, err = cache:get(constant.effectOn_key..id)
local duration, err = cache:get(constant.duration_key..id)
print("========eeo ", effectOn, " === ed ", duration)
if (not effectOn) or (not duration) then
    ngx.exit(ngx.HTTP_NOT_FOUND)
    print("exit caused by ineffectOn or duration")
    return
else 
    local now = ngx.now() * 1000
    local offset = effectOn - now
    if offset <= config.saltOffset and now < effectOn + duration then
        local cache = ngx.shared.scache
        local salt, err = cache:get(constant.salt_key..id)
        print("refreshing salt ", salt)
        if not salt then
            local uuid = require "uuid4"
            salt = uuid.getUUID()
            local success, err, _ = cache:set(constant.salt_key..id, salt)
            if not success then
                ngx.log(ngx.ERR, "can't set salt ", err)
            end
            print("setting salt ", salt)
        end
        salt, err = cache:get(constant.salt_key..id)
        print("set salt ", salt)
        ngx.say(string.format('{"effectOn":%d, "time":%d, "salt": "%s"}', effectOn, now, salt))
    else
        ngx.say(string.format('{"effectOn":%d, "time":%d}', effectOn, now))
    end
end

local args = ngx.req.get_uri_args()
-- parameters validations
function validate()
    if not args.id then
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
    end

    if not args.phone or not validatePhone(args.phone) then
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
    end

    if not args.salt then
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
    end
    local cache = ngx.shared.scache
    local salt, err = cache:get("salt:"..args.id)
    if not salt or args.salt ~= salt then
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
        return false
    end

    local effectOn, err = cache:get("eeo:"..args.id)
    if not effectOn then
        ngx.log(ngx.ERR, "can't get eeo ", err)
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
        return false
    end

    local duration, err = cache:get("ed:"..args.id)
    if not duration then
        ngx.log(ngx.ERR, "can't get ed ", err)
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
        return false
    end

    local now = ngx.now() * 1000
    if effectOn > now or now > effectOn + duration then 
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
        return false
    end
    return true
end

function validatePhone(phone)
    if string.sub(phone, 1, 2) == "86" then
        phone = string.sub(phone, 3)
    end
    if string.sub(phone, 1, 3) == "086" then
        phone = string.sub(phone, 4)
    end
    if string.match(phone, "^1[3|5|7|8|4]%d%d%d%d%d%d%d%d%d$") then
        return true
    else 
        return false
    end
end 

-- request limit counter update
function applyCounter()
    local counter = require "sk_counter"
    counter.apply(args.id)

    local count, _ = counter.get(args.id)
    ngx.log(ngx.INFO, "apply counter", count)
    if counter.stopped(args.id) then
        ngx.log(ngx.INFO, "stopped")
        ngx.exit(ngx.HTTP_FORBIDDEN)
    else 
        setToken()
    end
end

function setToken()
    local config = require "config"
    local ck = require "resty.cookie"
    local cookie, err = ck:new()
    if not cookie then
        ngx.log(ngx.CRIT, "can't create cookie ", err)
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
        return
    end

    local token = cookie:get(config.tokenCookie)
    if not token then
        local uuid = require "uuid4"
        local token = uuid.getUUID()
        local ok, err = cookie:set({key = config.tokenCookie, value = token, path = "/"})
        if not ok then
            ngx.log(ngx.CRIT, "can't set cookie ", err)
            ngx.exit(ngx.HTTP_NOT_ALLOWED)
            return
        end
    end
    setTokenStatus(token)
end

function setTokenStatus(token)
    local redisc = require "redisc"
    local redis = redisc:new()
    local ok, err = redis:hset("tk:"..token, "status", 1)
    if not ok then
        ngx.log(ngx.CRIT, "can't set token to redis ", err)
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
    end
end

function serve()
    if validate() then
        applyCounter()
    end
end

serve()

local util = require "access_util"
local args = ngx.req.get_uri_args()
-- parameters validations
function validate()
    if not args.id then
        ngx.log(ngx.INFO, "invalid id")
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
        return false
    end

    if not util.validatePhone(args.phone) then
        ngx.log(ngx.INFO, "invalid phone")
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
        return false
    end

    if not util.validateSalt(args.id, args.salt) then
        ngx.log(ngx.INFO, "invalid salt")
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
        return false
    end

    if not util.validateEffect(args.id) then
        ngx.log(ngx.INFO, "invalid effectOn")
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
        return false
    end
    return true
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
    ngx.log(ngx.INFO, "cookie ", token)

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
    ngx.log(ngx.INFO, "token set ", ok)
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

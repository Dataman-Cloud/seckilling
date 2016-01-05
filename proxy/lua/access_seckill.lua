-- parameters validations
function validate()
    local args = ngx.req.get_uri_args()
    if not args.id then
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
    end
    local cache = ngx.shared.scache
    local effectOn, err = cache:get("eeo:"..args.id)
    if not effectOn then
        ngx.log(ngx.ERR, "can't get eeo ", err)
        ngx.exit(ngx.HTTP_NOT_FOUND)
    end
    local duration, err = cache:get("ed:"..args.id)
    if not duration then
        ngx.log(ngx.ERR, "can't get ed ", err)
        ngx.exit(ngx.HTTP_NOT_FOUND)
    end
    local now = ngx.now() 
    print("======e ", effectOn, " d ", duration, "n ", now)
    if tonumber(effectOn) > now or now > tonumber(effectOn) + tonumber(duration) then 
        ngx.exit(ngx.HTTP_NOT_ALLOWED)
    end
    return true
end

-- request limit counter update
function applyCounter()
    local counter = require "sk_counter"
    counter.apply()

    local count, _ = counter.get()
    ngx.log(ngx.INFO, "apply counter", count)
    if counter.stopped() then
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
        ngx.exit(ngx.HTTP_NOT_ACCEPTABLE)
        return
    end

    local token = cookie:get(config.tokenCookie)
    if not token then
        local uuid = require "uuid4"
        local token = uuid.getUUID()
        local ok, err = cookie:set({key = config.tokenCookie, value = token, path = "/"})
        if not ok then
            ngx.log(ngx.CRIT, "can't set cookie ", err)
            ngx.exit(ngx.HTTP_NOT_ACCEPTABLE)
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
        ngx.exit(ngx.HTTP_NO_CONTENT)
    end
end

function serve()
    if validate() then
        applyCounter()
    end
end

serve()

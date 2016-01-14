local _M = {}
local config = require "config"
local constant = require "constant"


function _M.incr(eid)
    local cache = ngx.shared.scache
    local val, err = cache:incr(constant.count_key..eid, 1)
    if not val then
        ngx.log(ngx.ERR, "can't get counter", err)
    end

    return val, err
end

function _M.apply(eid) 
    local cache = ngx.shared.scache
    local val, err = cache:get(constant.count_key..eid)
    if not val then
        ngx.log(ngx.ERR, "can't get cache counter", err)
    end

    if val >= config.counterBatch then
        local redisc = require "redisc"
        local redis = redisc:new()

        local maxCount, err = redis:zcard(constant.sortset_key..eid)
        if not maxCount then
            ngx.log(ngx.CRIT, "can't get max_count ", err)
            ngx.exit(ngx.HTTP_NOT_ALLOWED)
            return
        end
        maxCount = tonumber(maxCount)

        local count, err = redis:incrby(constant.counter_key..eid, config.counterBatch)
        if not count then
            ngx.log(ngx.CRIT, "can't incrby redis", err)
        end
        ngx.log(ngx.INFO, "redis counter ", count, " maxCount ", maxCount)

        if count >= maxCount then
            local success, err, forcible = cache:set(constant.stop_key..eid, 1)
            if not success then
                ngx.log(ngx.ERR, "can't set stopped", err)
            end
        end
        _M.reset(eid)
    end
end

function _M.stopped(eid)
    local cache = ngx.shared.scache
    local val, err = cache:get(constant.stop_key..eid)
    if not val then
        ngx.log(ngx.CRIT, "can't get stopped", err)
    end

    return val == 1
end

function _M.enable(eid)
    local cache = ngx.shared.scache

    local success, err, forcible = cache:set(constant.stop_key..eid, 0)
    if not success then
        ngx.log(ngx.ERR, "can't clear stopped", err)
    end
end

function _M.reset(eid)
    local cache = ngx.shared.scache
    local success, err, forcible = cache:set(constant.count_key..eid, 0)
    if not success then
        ngx.log(ngx.ERR, "can't reset counter", err)
    end
end

function _M.get(eid)
    local cache = ngx.shared.scache
    local val, flags =  cache:get(constant.count_key..eid)
    return val
end

return _M

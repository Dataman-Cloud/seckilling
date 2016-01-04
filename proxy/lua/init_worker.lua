-- init worker
-- events, counter

local config = require "config"

function setEvents(redis)
    local testData = require "test_data"
    events = testData.events
    return redis:set("events", events)
end

function loadEvents()
    local mode = os.getenv("PROXY_MODE") 
    if not mode then
        mode = "prod"
    end
    ngx.log(ngx.INFO, "running in mode: ", mode)
    if mode == "dev" or mode == "DEV" then
        local testData = require "test_data"
        events = testData.events
    else
        local redisc = require "redisc"
        local redis = redisc:new()
        local events, err = redis:get("events")
        if not events then
            ngx.log(ngx.WARN, "can't retreive events from redis ", err)
        end
    end

    local cjson = require "cjson"
    local json = cjson.encode(events)

    ngx.log(ngx.INFO, "generated events json:\n\t", json)

    return json
end

function initEvents()
    local cache = ngx.shared.scache

    local val, err = cache:get("events")
    if not val then
        local json = loadEvents()
        cache:set("events", json)
        ngx.log(ngx.INFO, "set events successfully")
    end
end

function initCounter()
    local counter = require "sk_counter"
    counter.reset() 
    counter.enable()

    local redisc = require "redisc"
    local redis = redisc:new()

    local val, err = redis:set("counter", 0)
    if not val then
        ngx.log(ngx.CRIT, "can't reset redis counter ", err)
    end

    val, err = redis:set("max_count", config.maxCount)
    if not val then
        ngx.log(ngx.CRIT, "can't reset redis max_count ", err)
    end
    ngx.log(ngx.INFO, "counter reset")
end

function init()
    local cache = ngx.shared.scache
    cache:flush_all()

    ngx.log(ngx.INFO, "initializing server state...")

    local ok, err = ngx.timer.at(0,initEvents)
    if not ok then
        ngx.log(ngx.CRIT, "can't init events ", err)
    end
    ok, err = ngx.timer.at(0,initCounter)
    if not ok then
        ngx.log(ngx.CRIT, "can't init counter ", err)
    end

    ngx.log(ngx.INFO, "server state was initialized.")
end

init()

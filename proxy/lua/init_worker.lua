-- init worker
-- events, counter

function getEvents(redis)
    return redis:get("events")
end

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

function init()
    local cache = ngx.shared.scache
    cache:flush_all()
    ngx.log(ngx.INFO, "initializing server state...")
    local ok, err = ngx.timer.at(0,initEvents)
    ngx.log(ngx.INFO, "server state was initialized.")
end

init()

-- init worker
-- events, counter

local redisc = require "redisc"
local cjson = require "cjson"
local config = require "config"

function setEvents(redis)
    local testData = require "test_data"
    local events = testData.events
    return redis:set("events", events)
end

function loadEvents()
    local mode = os.getenv("PROXY_MODE") or "prod"
    ngx.log(ngx.INFO, "running in mode: ", mode)

    if mode == "dev" or mode == "DEV" then
        local testData = require "test_data"
        return testData.events
    else
        local redis = redisc:new()
        local eids, err = redis:lrange("events", 0, -1)
        if not eids then
            ngx.log(ngx.WARN, "can't retreive events from redis ", err)
        end
        return assambleEvents(redis, eids)
    end
end

function assambleEvents(redis, ids)
    local events = {}
    for i, id in ipairs(ids) do
        local res, err = redis:hgetall("event:" .. id)
        if not res then
            ngx.log(ngx.CRIT, "can't get event from redis id: ", id, " err: ", err)
        else 
            local redis = redisc:new()
            local event = redis:array_to_hash(res)
            event.id = tonumber(event.id)
            event.effectOn = tonumber(event.effectOn) * 1000
            event.duration = tonumber(event.duration) * 1000
            event.status = nil
            events[i] = event
        end
    end
    return events
end

function initEvents()
    local cache = ngx.shared.scache
    local val, err = cache:get("events")
    if not val then
        local events = loadEvents()
        local json = cjson.encode(events)
        ngx.log(ngx.INFO, "generated events json:\n\t", json)
        cache:set("events", json)
        setEventCache(cache, events)
        ngx.log(ngx.INFO, "set events successfully")
    end
end

function setEventCache(cache, events)
    local counter = require "sk_counter"
    for i = 1, #events do
        local event = events[i]
        cache:set("eeo:"..event.id, event.effectOn)
        cache:set("ed:"..event.id, event.duration)
        cache:set("count:"..event.id, 0)
        cache:set("stopped:"..event.id, 0)
        counter.reset(event.id)
        counter.enable(event.id)
    end
end

function init()
    local cache = ngx.shared.scache
    cache:flush_all()

    ngx.log(ngx.INFO, "initializing server state...")

    local ok, err = ngx.timer.at(0,initEvents)
    if not ok then
        ngx.log(ngx.CRIT, "can't init events ", err)
    end

    ngx.log(ngx.INFO, "server state was initialized.")
end

init()

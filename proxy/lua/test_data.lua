local _M = {}

local events = {}
local now = ngx.now() * 1000
events[1] =  {
    ["id"] = 1, 
    ["effectOn"] = now, 
    ["duration"] = 1000 * 60 * 8, 
    ["desc"] = "event 1", 
    ["resources"] = {
        ["small_image"] ="a1.png", 
        ["big_image"] = "a2.png"
    }
}

events[2] =  {
    ["id"] = 2, 
    ["effectOn"] = now + 15 * 1000 * 60, 
    ["duration"] = 1000 * 60 * 8, 
    ["desc"] = "event 2", 
    ["resources"] = {
        ["small_image"] ="b1.png", 
        ["big_image"] = "b2.png"
    }
}

events[2] =  {
    ["id"] = 2, 
    ["effectOn"] = now + 15 * 1000 * 60 * 2,
    ["duration"] = 1000 * 60 * 8, 
    ["desc"] = "event 3", 
    ["resources"] = {
        ["small_image"] ="c1.png", 
        ["big_image"] = "c2.png"
    }
}

_M.events = events

return _M

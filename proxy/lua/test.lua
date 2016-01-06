local counter = require "sk_counter"
local cnt = counter.get(1)

ngx.say(cnt)

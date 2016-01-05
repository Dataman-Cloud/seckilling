local counter = require "sk_counter"
local cnt = counter.get()

ngx.say(cnt)

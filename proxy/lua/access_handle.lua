local counter = require "sk_counter"
counter.apply()
local count, _ = counter.get()
ngx.log(ngx.INFO, "apply counter", count)
if counter.stopped() then
    ngx.log(ngx.INFO, "stopped")
    ngx.exit(ngx.HTTP_FORBIDDEN)
end

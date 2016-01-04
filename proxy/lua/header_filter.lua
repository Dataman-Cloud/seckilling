local status = ngx.status

if status == ngx.HTTP_OK then
    local counter = require "sk_counter"
    local cnt, err = counter.incr()

    ngx.log(ngx.INFO, "=========counter incr: ", cnt)
end

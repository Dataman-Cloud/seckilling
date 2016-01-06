local status = ngx.status
local args = ngx.req.get_uri_args()

if status == ngx.HTTP_OK then
    local counter = require "sk_counter"
    local cnt, err = counter.incr(args.id)

    ngx.log(ngx.INFO, "=========counter incr: ", cnt)
end

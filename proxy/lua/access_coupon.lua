local util = require "access_util"
local args = ngx.req.get_uri_args()

if not args.id then
    ngx.log(ngx.INFO, "invalid id")
    ngx.exit(ngx.HTTP_NOT_ALLOWED)
    return false
end

if not util.validatePhone(args.phone) then
    ngx.exit(ngx.HTTP_NOT_ALLOWED)
    return
end


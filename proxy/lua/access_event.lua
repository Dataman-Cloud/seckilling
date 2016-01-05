local args = ngx.req.get_uri_args()
if not args.id then
    ngx.exit(ngx.HTTP_NOT_ALLOWED)
end

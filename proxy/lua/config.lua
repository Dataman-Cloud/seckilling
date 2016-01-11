local _M = {}

-- redis connection
local redis = {} 
-- redis.host = "123.59.58.58"
redis.host = os.getenv("REDIS_HOST") or "123.59.58.58"
--redis.port = 5506
redis.port =  os.getenv("REDIS_PORT") or 5506
--redis.password = "UQPqcj7nUyii38cpYcr9OnTbIJ3dHXvJ"
_M.redis = redis

_M.counterBatch = 5
_M.maxCount = 10
_M.saltOffset = 1000 * 60 * 3
_M.tokenCookie = 'DM_SK_UID'

return _M

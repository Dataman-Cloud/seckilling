local _M = {}

-- redis connection
local redis = {} 
redis.host = "123.59.58.58"
redis.port = 5506
redis.password = "UQPqcj7nUyii38cpYcr9OnTbIJ3dHXvJ"
_M.redis = redis

_M.counterBatch = 5
_M.maxCount = 10

return _M

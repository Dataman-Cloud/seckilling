local _M = {}

-- redis connection
local redis = {} 
-- redis.host = "123.59.58.58"
redis.host = os.getenv("REDIS_HOST") or "123.59.58.58"
--redis.port = 5506
redis.port =  os.getenv("REDIS_PORT") or 5506
--redis.password = "UQPqcj7nUyii38cpYcr9OnTbIJ3dHXvJ"
_M.redis = redis

local counterBatch = os.getenv("COUNTER_BATCH")
_M.counterBatch = tonumber(counterBatch) or 5

local saltOffset = os.getenv("SALT_OFFSET")
_M.saltOffset = tonumber(saltOffset) or 180000 --1000 * 60 * 3
_M.tokenCookie = os.getenv("TOKEN_COOKIE") or "DM_SK_UID"

return _M

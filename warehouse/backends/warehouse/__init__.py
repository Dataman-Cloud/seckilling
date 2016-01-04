import redis

from django.conf import settings

redis_pool = redis.ConnectionPool(
    host=settings.REDIS['host'],
    port=settings.REDIS['port'],
    db=settings.REDIS['db']
)

redis_inst = redis.Redis(connection_pool=redis_pool)

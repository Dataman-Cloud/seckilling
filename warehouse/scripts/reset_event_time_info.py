import redis
import time

r = redis.StrictRedis(host='123.59.58.58', port=5506)

START_AHEAD = 10*60      # n minutes ahead the current time for the start of the first event
DURATION = 4*60          # n minutes for each round of event

# reset event hash time info
start_timestamp = int(time.time()) + START_AHEAD
for idx, event_id in enumerate(r.lrange('events', 0, -1)):
    mapping = {
        'effectOn': start_timestamp + idx * DURATION,
        'duration': DURATION
    }
    r.hmset('event:' + event_id.decode(encoding='UTF-8'), mapping=mapping)

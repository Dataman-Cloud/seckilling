import itertools
import logging
import time

from django.conf import settings
from django.db.models.signals import post_save, post_delete
from django.dispatch import receiver
from django.db import models, transaction
from django.utils import timezone
from django.db.models import Sum, Q
from django.core.exceptions import ValidationError

from . import redis_inst


logger = logging.getLogger(__name__)

# Create your models here.
PRIZE_LEVEL = [
    (0, 'none'),
    (1, 'first'),
    (2, 'second'),
    (3, 'third')
]


class Brand(models.Model):
    brand_id = models.CharField(max_length=128)  #渠道ID
    name = models.CharField(max_length=128)  #渠道名称
    logo = models.CharField(max_length=128)  #渠道logo链接
    exchange_link = models.CharField(max_length=128)  #渠道兑奖链接
    exchange_detail = models.TextField()  #渠道兑奖须知

    @property
    def delivered_prizes(self):
        return self.prizes.filter(is_taken=True)

    def __str__(self):
        return self.name


class Activities(models.Model):
    class Meta:
        ordering = ["start_at"]

    STATUS = [
        ("waiting", 'waiting'),
        ("running", 'running'),
        ("end", 'end')
    ]

    name = models.CharField(max_length=128)
    start_at = models.DateTimeField(null=False)
    end_at = models.DateTimeField(null=False)
    status = models.CharField(max_length=255, choices=STATUS, default='waiting')

    def __str__(self):
        return ';'.join([str(self.id), self.name, self.status])

    def clean(self):
        if self.start_at > self.end_at:
            raise ValidationError('The end time should be later than the start time')


class Activities_item(models.Model):

    brand = models.ForeignKey(Brand)
    level = models.IntegerField(choices=PRIZE_LEVEL, default=0)
    count = models.PositiveIntegerField(default=0, null=False)
    activity = models.ForeignKey(Activities, related_name='items')

    def __str__(self):
        return ';'.join([str(self.id), self.brand.name, str(self.count)])

    def clean(self):
        if not hasattr(self, 'brand'):
            raise ValidationError('Need choose the Brand!')

        prizes = Prizes.objects.filter(brand=self.brand).filter(level=self.level)
        items = Activities_item.objects.filter(brand=self.brand).filter(level=self.level)

        if prizes:
            prizes_total = prizes.count()

            if items:
                prizes_taken = items.aggregate(Sum('count'))['count__sum']
            else:
                prizes_taken = 0

            if self.id:
                this_item = Activities_item.objects.get(id=self.id)
                prizes_taken = prizes_taken - this_item.count
                prizes_available = prizes_total - prizes_taken
            else:
                prizes_available = prizes_total - prizes_taken

            count_sum = prizes_taken + self.count

            if count_sum > prizes_total:
                raise ValidationError('There is not enough prizes, only %d available in %d/%d' %
                                      (prizes_available, prizes_taken, prizes_total))
        else:
            raise ValidationError('wrong conditions, brand: %s, level: %s' % (self.brand.name, str(self.level)))

    @property
    def delivered_prize_count(self):
        # NOTE(xychu): If we use pub/sub to send msg, then no need to check out the status.
        if self.status == 'end':
            return self.prizes.filter(~Q(winner_cell='')).count()
        else:
            # NOTE(xchu): Not started yet or not finished, we only collect result *after* event finished;
            #             For running event, we have a API to be called to get current count in Redis.
            return 0


class Prizes(models.Model):

    class Meta:
        unique_together = ('brand', 'exchange_code')

    prize_id = models.CharField(max_length=128, unique=True)  # 奖品ID
    name = models.CharField(max_length=128)  # 奖品名字
    exchange_code = models.CharField(max_length=128, null=False)  # 奖品兑换码
    thumbnail_path = models.CharField(max_length=128)  # 奖品展示图
    detail = models.TextField()  # 奖品展示图
    brand = models.ForeignKey(Brand, related_name='prizes')  # 兑奖渠道
    winner_cell = models.CharField(max_length=20, blank=True)  # 中奖手机号
    level = models.IntegerField(choices=PRIZE_LEVEL, default=0)
    value = models.PositiveIntegerField(null=True, blank=True)  # 奖品金额
    created_at = models.DateTimeField(default=timezone.now, null=True)
    is_taken = models.BooleanField(default=False, null=False)
    taken_at = models.DateTimeField(blank=True, null=True)
    activity_item = models.ForeignKey(Activities_item, null=True, on_delete=models.SET_NULL,
                                      related_name='prizes', blank=True)

    def __str__(self):
        return '{0}: {1}'.format(self.brand.name, self.exchange_code)


@receiver(post_save, sender=Activities_item)
def update_prizes(sender, instance, created, **kwargs):
    old_prize_ids = None
    logger.info("Post save signal on event[%s] item[%s] received." % (instance.activity.id, instance.id))
    if created:
        logger.info('Event[%s] item[%s] created, arranging...' % (instance.activity.id, instance.id))
    else:
        logger.info('Event[%s] item[%s] updated, rearranging...' % (instance.activity.id, instance.id))
        old_prize_ids = list(instance.prizes.values_list('prize_id', flat=True))
        instance.prizes.update(activity_item=None)

    id_qs = Prizes.objects.filter(brand=instance.brand,
                                  level=instance.level,
                                  activity_item__isnull=True).values_list('id', flat=True)
    id_list = list(id_qs)[:instance.count]
    qs = Prizes.objects.filter(id__in=id_list)
    qs.update(activity_item=instance)

    logger.info("Assign %s SNs for event[%s] item[%s] done." % (instance.count, instance.activity.id, instance.id))

    # load SNs for this event item
    sn_key = settings.REDIS['key_fmts']['sn_set'] % str(instance.activity.id)
    if old_prize_ids:
        logger.info('Start removing %s SNs from from redis.' % len(old_prize_ids))
        redis_inst.zrem(sn_key, *old_prize_ids)
        logger.info('Removing %s SNs from from redis done.' % len(old_prize_ids))
    source_list = list(itertools.chain.from_iterable([[prize.prize_id,
                                                       time.time()] for prize in instance.prizes.all()]))
    if source_list:
        redis_inst.zadd(sn_key, *source_list)
        logger.info("Load %s SNs into redis %s done." % (instance.prizes.count(), sn_key))


@receiver(post_save, sender=Activities)
def update_prize_and_load_to_redis(sender, instance, created, **kwargs):

    logger.info("Post save signal on event[%s] received." % instance)

    if kwargs['update_fields'] and 'status' in kwargs['update_fields'] and instance.status == 'end':
        logger.info("Event[%s] is ended. Start pulling results back...")
        # updating activity status, we should pull the results back *only if* it's ended.
        with transaction.atomic():
            valid_count = 0
            for item in instance.items.all():
                for prize_id, exchange_code in item.prizes.values_list('prize_id', 'exchange_code'):
                    winner_cell = redis_inst.hget(settings.REDIS['key_fmts']['result_hash'] % (instance.id, prize_id),
                                                  settings.REDIS['key_fmts']['cell_key'])
                    if not winner_cell:
                        logger.warn("No result[%s] found in redis [%s]." %
                                    (settings.REDIS['key_fmts']['cell_key'],
                                     settings.REDIS['key_fmts']['result_hash'] % (instance.id, prize_id)))
                    else:
                        winner_cell = winner_cell
                        item.prizes.filter(prize_id=prize_id).update(winner_cell=winner_cell)
                        valid_count += 1
            logger.info("%s results pulled back for event[%s]." % (valid_count, instance))
    else:
        # create/update activity before starting, we need to reload data into redis
        if created:
            logger.info('Event[%s] created, arranging...' % instance)

            # clean sn_key for this event item, values will be added in Activity_item post save.
            sn_key = settings.REDIS['key_fmts']['sn_set'] % str(instance.id)
            redis_inst.delete(sn_key)
        else:
            logger.info('Event[%s] updated, arranging...' % instance)
        # Load event data to redis
        # set event hash with key: event:<e_id>
        event_key = settings.REDIS['key_fmts']['event_hash'] % str(instance.id)
        redis_inst.delete(event_key)
        mapping = {
            'id': instance.id,
            'effectOn': int(instance.start_at.timestamp()),
            'duration': int(instance.end_at.timestamp()) - int(instance.start_at.timestamp()),
            'desc': ''
        }
        redis_inst.hmset(event_key, mapping=mapping)

        logger.info('Adding event hash[%s]: %s' % (event_key, mapping))

        # update events list
        load_events()

        if not created:
            logger.info('Event[%s] rearranging done.' % instance)
        else:
            logger.info('Event[%s] arranging done.' % instance)


@receiver(post_delete, sender=Activities)
def clean_up_redis(sender, instance, **kwargs):
    logger.info("Post delete signal on event[%s] received." % instance)
    # delete event hash
    event_key = settings.REDIS['key_fmts']['event_hash'] % str(instance.id)
    redis_inst.delete(event_key)
    logger.info("Event hash[%s] deleted." % event_key)

    # delete SNs for this event
    sn_key = settings.REDIS['key_fmts']['sn_set'] % str(instance.id)
    redis_inst.delete(sn_key)
    logger.info("SNs Sorted Sets[%s] deleted." % sn_key)

    # update events list
    load_events()


def load_events():
    events = Activities.objects.values_list('id', flat=True)
    events_key = settings.REDIS['key_fmts']['events_list']

    if not events:
        logger.info("No events left, reload will delete the events list in redis.")

    redis_inst.delete(events_key)
    if events:
        logger.info("Need to reload events list %s in redis %s. Reloading..." % (events, events_key))
        redis_inst.rpush(events_key, *events)
        logger.info("Reload events list %s in redis %s done." % (events, events_key))

import itertools

from django.conf import settings
from django.db.models.signals import post_save
from django.dispatch import receiver
from django.db import models
from django.utils import timezone
from django.db.models import Sum
from django.core.exceptions import ValidationError

from . import redis_inst


# Create your models here.
PRIZE_LEVEL = [
    (0, 'none'),
    (1, 'first'),
    (2, 'second'),
    (3, 'third')
]


class Brand(models.Model):
    name = models.CharField(max_length=128)

    @property
    def delivered_prizes(self):
        return self.prizes.filter(is_taken=True)

    def __str__(self):
        return self.name


class Activities(models.Model):
    STATUS = [
        ("waiting", 'waiting'),
        ("running", 'running'),
        ("end", 'end')
    ]

    class Meta:
        ordering = ["start_at"]

    start_at = models.DateTimeField(null=False)
    end_at = models.DateTimeField(null=False)
    brand = models.ForeignKey(Brand)
    level = models.IntegerField(choices=PRIZE_LEVEL, default=0)
    count = models.PositiveIntegerField(default=0, null=False)
    status = models.CharField(max_length=255, choices=STATUS, default='waiting')

    def __str__(self):
        return ';'.join([str(self.id), self.brand.name, str(self.count), self.status])

    def clean(self):
        if not hasattr(self, 'brand'):
            raise ValidationError('Need choose the Brand!')

        prizes = Prizes.objects.filter(brand=self.brand).filter(level=self.level)
        activates = Activities.objects.filter(brand=self.brand).filter(level=self.level)

        if prizes:
            prizes_total = prizes.count()

            if activates:
                prizes_taken = activates.aggregate(Sum('count'))['count__sum']
            else:
                prizes_taken = 0

            if self.id:
                this_activate = Activities.objects.get(id = self.id)
                prizes_taken = prizes_taken - this_activate.count
                prizes_avaliable = prizes_total - prizes_taken
            else:
                prizes_avaliable = prizes_total - prizes_taken

            count_sum = prizes_taken + self.count

            if self.start_at > self.end_at:
                raise ValidationError('the end time should be later than the start time')
            elif count_sum > prizes_total:
                raise ValidationError('There is not enough prizes, only %d avaliable in %d/%d' % (prizes_avaliable, prizes_taken, prizes_total))
        else:
            raise ValidationError('wrong conditions, brand: %s, level: %s' % (self.brand.name, str(self.level)))

    # def save(self, **kwargs):
    #     is_update = False
    #     if self.pk:
    #         is_update = True
    #     super(Activities, self).save(**kwargs)
    #     import ipdb;ipdb.set_trace()
    #     # set prize activity
    #     if is_update:
    #         self.prizes.update(activity=None)
    #
    #     # count check
    #     self.clean()
    #
    #     id_list = Prizes.objects.filter(brand=self.brand,
    #                                     level=self.level,
    #                                     activity__isnull=True).values_list('id')[:self.count]
    #     qs = Prizes.objects.filter(id__in=id_list)
    #     qs.update(activity=self)


class Prizes(models.Model):

    class Meta:
        unique_together = ('brand', 'serial_number')

    serial_number = models.CharField(max_length=128, null=False)
    brand = models.ForeignKey(Brand, related_name='prizes')
    level = models.IntegerField(choices=PRIZE_LEVEL, default=0)
    created_at = models.DateTimeField(default=timezone.now, null=True)
    is_taken = models.BooleanField(default=False, null=False)
    taken_at = models.DateTimeField(blank=True, null=True)
    winner_cell = models.CharField(max_length=20, blank=True)

    activity = models.ForeignKey(Activities, null=True, related_name='prizes')

    def __str__(self):
        return '{0}: {1}'.format(self.brand.name, self.serial_number)

@receiver(post_save, sender=Activities)
def update_prize_activity(sender, instance, **kwargs):
    # FIXME(xychu): need to support update
    id_qs = Prizes.objects.filter(brand=instance.brand,
                                  level=instance.level,
                                  activity__isnull=True).values_list('id', flat=True)
    id_list = list(id_qs)[:instance.count]
    qs = Prizes.objects.filter(id__in=id_list)
    qs.update(activity=instance)

@receiver(post_save, sender=Activities)
def load_to_redis(sender, instance, **kwargs):
    # set event hash with key: event:<e_id>
    event_key = 'event:' + str(instance.id)
    redis_inst.delete(event_key)
    mapping = {
        'id': instance.id,
        'effectOn': int(instance.start_at.timestamp()),
        'duration': int(instance.end_at.timestamp()) - int(instance.start_at.timestamp()),
        'desc': ''
    }
    redis_inst.hmset(event_key, mapping=mapping)

    # load SNs for this event
    sn_key = 'SN:' + str(instance.id)
    redis_inst.delete(sn_key)
    source_list = list(itertools.chain.from_iterable(
            [[item.serial_number, item.id] for item in instance.prizes.all()]))
    redis_inst.zadd(sn_key, *source_list)


def pull_result_back(key):
    raise NotImplementedError


def get_current_activity():
    return redis_inst.get(settings.REDIS['current_eid'])


def load_events():
    events = Activities.objects.values_list('id', flat=True)
    redis_inst.delete('events')
    redis_inst.rpush('events', *events)

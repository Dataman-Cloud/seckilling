from django.db import models
from django.utils import timezone
from django.db.models import Sum
from django.core.exceptions import ValidationError


# Create your models here.


class Brand(models.Model):
    name = models.CharField(max_length=128)


class Prizes(models.Model):

    class Meta:
        unique_together = ('brand', 'serial_number')

    PRIZE_LEVEL = [
        (0, 'none'),
        (1, 'first'),
        (2, 'second'),
        (3, 'third')
    ]

    serial_number = models.CharField(max_length=128, null=False)
    brand = models.ForeignKey(Brand, related_name='prizes')
    level = models.IntegerField(choices=PRIZE_LEVEL, default=0)
    created_at = models.DateTimeField(default=timezone.now, null=True)
    is_taken = models.BooleanField(default=False, null=False)
    taken_at = models.DateTimeField(blank=True, null=True)
    winner_cell = models.CharField(max_length=20, blank=True)

    def __str__(self):
        return ''.join([self.brand.name, '/', self.serial_number])


class Activities(models.Model):
    class Meta:
        ordering = ["start_at"]

    PRIZE_LEVEL = [
        (0, 'none'),
        (1, 'first'),
        (2, 'second'),
        (3, 'third')
    ]

    start_at = models.DateTimeField(unique=True, null=False)
    end_at = models.DateTimeField(unique=True, null=False)
    brand = models.ForeignKey(Brand)
    level = models.IntegerField(choices=PRIZE_LEVEL, default=0)
    count = models.PositiveIntegerField(default=0, null=False)

    def __str__(self):
        return ','.join([self.brand.name, str(self.count)])

    def clean(self):
        activates = Activities.objects.all()
        if activates:
            count_sum = Activities.objects.filter(brand=self.brand).aggregate(Sum('count'))
            count_sum = count_sum['count__sum'] + self.count
            print(count_sum)
            print(self.start_at, self.end_at)
            if self.start_at > self.end_at:
                raise ValidationError('the end time should be later than the start time')
            elif count_sum > Prizes.objects.filter(brand=self.brand).count():
                raise ValidationError('There is not enough prizes')

from django.db import models
from django.utils import timezone

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

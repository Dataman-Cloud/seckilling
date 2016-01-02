from django.db import models
from django.utils import timezone

# Create your models here.

import datetime

class Base(models.Model):

    class Meta:
       abstract = True
    

class Prizes(Base):

    class Meta:
       db_table = 'prizes'

    PRIZENAME = [ 
        ('special', 'special'), 
        ('first', 'first'),
        ('second', 'second'), 
        ('third', 'third')
     ]

    id = models.AutoField(primary_key=True)
    prize_code = models.CharField(max_length=128, unique=True, null=False)
    prize_brand = models.CharField(max_length=128, null=False)
    prize_name = models.CharField(max_length=128, choices=PRIZENAME, default='third', null=False) 
    is_taken = models.BooleanField(default=False, null=False) 
    created_at = models.DateTimeField(default=timezone.now, null=True)
    
class Winners(Base):

    class Meta:
       db_table = 'winners'
    
    id = models.AutoField(primary_key=True)
    phone_number = models.CharField(max_length=20, unique=True, null=True)
    prize_time = models.DateTimeField(default=timezone.now, null=False)
    prize_code =  models.ForeignKey(Prizes)

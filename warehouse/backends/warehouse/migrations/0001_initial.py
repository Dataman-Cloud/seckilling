# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import migrations, models
import django.utils.timezone


class Migration(migrations.Migration):

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Prizes',
            fields=[
                ('id', models.AutoField(primary_key=True, serialize=False)),
                ('prize_code', models.CharField(max_length=128, unique=True)),
                ('prize_brand', models.CharField(max_length=128)),
                ('prize_name', models.CharField(max_length=128, choices=[('special', 'special'), ('first', 'first'), ('second', 'second'), ('third', 'third')], default='third')),
                ('is_taken', models.BooleanField()),
            ],
            options={
                'db_table': 'prizes',
            },
        ),
        migrations.CreateModel(
            name='Winners',
            fields=[
                ('id', models.AutoField(primary_key=True, serialize=False)),
                ('phone_number', models.CharField(null=True, max_length=20, unique=True)),
                ('prize_time', models.DateTimeField(default=django.utils.timezone.now)),
                ('prize_code', models.ForeignKey(to='warehouse.Prizes')),
            ],
            options={
                'db_table': 'winners',
            },
        ),
    ]

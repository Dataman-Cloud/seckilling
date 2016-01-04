# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('warehouse', '0002_auto_20160102_1226'),
    ]

    operations = [
        migrations.CreateModel(
            name='Activities',
            fields=[
                ('id', models.AutoField(serialize=False, verbose_name='ID', primary_key=True, auto_created=True)),
                ('start_at', models.DateTimeField()),
                ('end_at', models.DateTimeField()),
                ('level', models.IntegerField(choices=[(0, 'none'), (1, 'first'), (2, 'second'), (3, 'third')], default=0)),
                ('count', models.PositiveIntegerField(default=0)),
                ('status', models.CharField(max_length=255, choices=[('waiting', 'waiting'), ('running', 'running'), ('end', 'end')], default='waiting')),
                ('brand', models.ForeignKey(to='warehouse.Brand')),
            ],
            options={
                'ordering': ['start_at'],
            },
        ),
    ]

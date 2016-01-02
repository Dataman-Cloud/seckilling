# -*- coding: utf-8 -*-
from __future__ import unicode_literals

from django.db import migrations, models
import django.utils.timezone


class Migration(migrations.Migration):

    dependencies = [
        ('warehouse', '0001_initial'),
    ]

    operations = [
        migrations.AddField(
            model_name='prizes',
            name='created_at',
            field=models.DateTimeField(null=True, default=django.utils.timezone.now),
        ),
        migrations.AlterField(
            model_name='prizes',
            name='is_taken',
            field=models.BooleanField(default=False),
        ),
    ]

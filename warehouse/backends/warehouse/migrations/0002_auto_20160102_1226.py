# -*- coding: utf-8 -*-
# Generated by Django 1.9 on 2016-01-02 12:26
from __future__ import unicode_literals

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('warehouse', '0001_initial'),
    ]

    operations = [
        migrations.AlterField(
            model_name='prizes',
            name='taken_at',
            field=models.DateTimeField(blank=True, null=True),
        ),
    ]
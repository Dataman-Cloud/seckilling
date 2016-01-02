# -*- coding: utf-8 -*-
# Generated by Django 1.9 on 2016-01-02 12:22
from __future__ import unicode_literals

from django.db import migrations, models
import django.db.models.deletion
import django.utils.timezone


class Migration(migrations.Migration):

    initial = True

    dependencies = [
    ]

    operations = [
        migrations.CreateModel(
            name='Brand',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('name', models.CharField(max_length=128)),
            ],
        ),
        migrations.CreateModel(
            name='Prizes',
            fields=[
                ('id', models.AutoField(auto_created=True, primary_key=True, serialize=False, verbose_name='ID')),
                ('serial_number', models.CharField(max_length=128)),
                ('level', models.IntegerField(choices=[(0, 'none'), (1, 'first'), (2, 'second'), (3, 'third')], default=0)),
                ('created_at', models.DateTimeField(default=django.utils.timezone.now, null=True)),
                ('is_taken', models.BooleanField(default=False)),
                ('taken_at', models.DateTimeField(blank=True)),
                ('winner_cell', models.CharField(blank=True, max_length=20)),
                ('brand', models.ForeignKey(on_delete=django.db.models.deletion.CASCADE, related_name='prizes', to='warehouse.Brand')),
            ],
        ),
        migrations.AlterUniqueTogether(
            name='prizes',
            unique_together=set([('brand', 'serial_number')]),
        ),
    ]

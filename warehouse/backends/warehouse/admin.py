from django.contrib import admin

# Register your models here.
from .models import Prizes, Brand, Activities


admin.site.register(Brand)
admin.site.register(Prizes)
admin.site.register(Activities)

from django.contrib import admin

# Register your models here.
from .models import Prizes, Brand


admin.site.register(Brand)
admin.site.register(Prizes)



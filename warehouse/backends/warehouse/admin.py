from django.contrib import admin

# Register your models here.
from .models import Prizes
from .models import Winners

admin.site.register(Prizes)

admin.site.register(Winners)



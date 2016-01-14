from django.contrib import admin

# Register your models here.
from .models import Prizes, Brand, Activities, Activities_item

class ActivitiesItemInline(admin.TabularInline):
    model = Activities_item
    extra = 0

class PrizesAdmin(admin.ModelAdmin):

    list_display = ('prize_id', 'name', 'exchange_code',
                    'brand', 'level', 'is_taken',
                    'winner_cell', 'activity)

    list_filter = ['name', 'brand', 'is_taken', 'level', 'activity']
    search_fields = ['exchange_code', 'winner_cell']


class ActivitiesAdmin(admin.ModelAdmin):

    list_display = ('name', 'start_at', 'end_at', 'status'
                   )

    search_fields = ['name']
    list_filter = ['status']

    inlines = [ActivitiesItemInline]

class BrandAdmin(admin.ModelAdmin):

    list_display = ('brand_id', 'name', )

    search_fields = ['brand_id', 'name']


admin.site.register(Prizes, PrizesAdmin)
admin.site.register(Brand, BrandAdmin)
admin.site.register(Activities, ActivitiesAdmin)

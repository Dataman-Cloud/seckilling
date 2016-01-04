from django.contrib import admin

# Register your models here.
from .models import Prizes, Brand, Activities

admin.site.register(Brand)

class PrizesAdmin(admin.ModelAdmin):

    readonly_fields = ('serial_number', 'brand', 'level', 'created_at',
                       'is_taken', 'taken_at', 'winner_cell', 'activity'
                      )

    list_display = ('id', 'brand', 'serial_number', 'level', 'is_taken',
                    'winner_cell', 'activity')

    list_filter = ['brand', 'is_taken', 'level', 'activity']

    search_fields = ['serial_number', 'winner_cell']


class ActivitiesAdmin(admin.ModelAdmin):

    list_display = ('id', 'start_at', 'end_at', 'brand', 'level', 'count',
                     'status'
                   )

    list_filter = ['brand', 'level', 'status']


admin.site.register(Prizes, PrizesAdmin)
admin.site.register(Activities, ActivitiesAdmin)

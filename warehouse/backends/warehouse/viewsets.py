from django.db.models import Count

from django.conf import settings
from rest_framework import viewsets
from rest_framework.decorators import detail_route, list_route
from rest_framework.response import Response

from .serializer import BrandStatsSerializer, PrizeSerializer, ActivitiesSerializer
from .models import Brand, Prizes, Activities

from . import redis_inst

class BrandStatsViewSet(viewsets.ReadOnlyModelViewSet):
    """
    """
    # def list(self, request, *args, **kwargs):
    #     pass
    #
    # def retrieve(self, request, *args, **kwargs):
    #     pass

    def get_queryset(self):
        return Brand.objects.annotate(total_prize_count=Count('prizes'))

    queryset = Brand.objects.all()
    serializer_class = BrandStatsSerializer


class PrizeViewSet(viewsets.ReadOnlyModelViewSet):
    """
    """
    # def list(self, request, *args, **kwargs):
    #     pass
    #
    # def retrieve(self, request, *args, **kwargs):
    #     pass

    queryset = Prizes.objects.all()
    serializer_class = PrizeSerializer

class ActivitiesViewSet(viewsets.ReadOnlyModelViewSet):
    """
    """
    # def list(self, request, *args, **kwargs):
    #     pass
    #
    # def retrieve(self, request, *args, **kwargs):
    #     pass

    @detail_route(methods=['get'], url_path='delivered-count')
    def delivered_count(self, request, pk, **kwages):
        event_key = settings.REDIS['key_fmts']['delivered_count'] % str(pk)
        count = redis_inst.get(event_key)
        return Response({'count': int(count)})

    @list_route(methods=['get'], url_path='current-activity')
    def current_activity(self, request, **kwages):
        event_key = settings.REDIS['key_fmts']['current_eid']
        current_eid = int(redis_inst.get(event_key))
        event_key = settings.REDIS['key_fmts']['event_hash'] % str(int(current_eid))
        status = int(redis_inst.hget(event_key, 'status'))
        update_activity_status(current_eid, status)
        return Response({'current_eid': current_eid})

    queryset = Activities.objects.all()
    serializer_class = ActivitiesSerializer


def update_activity_status(eid, status, **kwargs):
    event_status_mapping = {
        1: 'running',
        2: 'waiting',
        3: 'end'
    }
    qs = Activities.objects.get(id=eid)
    status = event_status_mapping[status]
    if qs.status != status:
        qs.status=status
        qs.save()

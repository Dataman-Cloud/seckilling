from django.db.models import Count

from rest_framework import viewsets

from .serializer import BrandStatsSerializer, PrizeSerializer
from .models import Brand, Prizes


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

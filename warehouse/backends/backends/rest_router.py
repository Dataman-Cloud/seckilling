from rest_framework import routers

from warehouse import viewsets as wh_viewsets

router = routers.DefaultRouter(trailing_slash=False)

router.register(r'brand-stats', wh_viewsets.BrandStatsViewSet)
router.register(r'prizes', wh_viewsets.PrizeViewSet)
router.register(r'activities', wh_viewsets.ActivitiesViewSet)

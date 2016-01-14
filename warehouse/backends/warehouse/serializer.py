from rest_framework import serializers

from .models import Prizes, Brand, Activities


class BrandStatsSerializer(serializers.ModelSerializer):
    total_prize_count = serializers.IntegerField()
    delivered_prize_count = serializers.SerializerMethodField()

    def get_delivered_prize_count(self, obj):
        return obj.prizes.filter(is_taken=True).count()

    class Meta:
        model = Brand
        fields = ('name', 'total_prize_count', 'delivered_prize_count')


class PrizeSerializer(serializers.ModelSerializer):
    brand = serializers.CharField(source='brand.name')

    class Meta:
        model = Prizes
        fields = ('exchange_code', 'brand', 'level', 'winner_cell', 'created_at')


class ActivitiesSerializer(serializers.ModelSerializer):
    brand = serializers.CharField(source='brand.name')
    delivered_prize_count = serializers.IntegerField()

    class Meta:
        model = Activities
        fields = ('id', 'start_at', 'end_at', 'brand', 'level', 'count', 'status', 'delivered_prize_count')

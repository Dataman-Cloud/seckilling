from rest_framework import serializers

from .models import Prizes, Brand


class BrandStatsSerializer(serializers.ModelSerializer):
    total_prize_count = serializers.SerializerMethodField()
    delivered_prize_count = serializers.SerializerMethodField()

    def get_total_prize_count(self, obj):
        return obj.prizes.count()

    def get_delivered_prize_count(self, obj):
        return obj.prizes.filter(is_taken=True).count()

    class Meta:
        model = Brand
        fields = ('name', 'total_prize_count', 'delivered_prize_count')


class PrizeSerializer(serializers.ModelSerializer):
    brand = serializers.CharField(source='brand.name')

    class Meta:
        model = Prizes
        fields = ('serial_number', 'brand', 'level', 'winner_cell', 'created_at')

import uuid
import random

from django.http import HttpResponse

from .models import Prizes, Brand
# Create your views here.


def index(request):
    return HttpResponse("Hello")


def gen_data(request):
    """
    Test only.
    """
    target_count = 300000
    current_count = Prizes.objects.count()
    if current_count < target_count:
        for brand in ["meituan", "baidu", "tmall"]:
            Brand.objects.get_or_create(name=brand)
        prizes = []
        for i in range(target_count - current_count):
            sn = uuid.uuid4().hex
            brand = random.choice(list(Brand.objects.all()))
            prizes.append(Prizes(serial_number=sn, brand=brand))
        try:
            Prizes.objects.bulk_create(prizes)
        except Exception as e:
            return HttpResponse(e, status=500)
        else:
            return HttpResponse("Gen data", status=201)
    else:
        return HttpResponse("Already have enough data.")

import uuid
import random

from django.shortcuts import render

from django.http import HttpResponse

from .models import Prizes
# Create your views here.

def index(request):
    return HttpResponse("Hello")

def gendata(request):
    for i in range(1, 2):
        prize_code = uuid.uuid4().hex
        prize_brand = random.choice(("meituan", "baidu", "koubei"))
        prize = Prizes.objects.create(prize_code=prize_code, prize_brand=prize_brand)
    return HttpResponse("Gen data")
    

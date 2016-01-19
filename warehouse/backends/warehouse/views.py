import uuid
import random

from django.shortcuts import render, render_to_response

from django.http import HttpResponse, HttpResponseRedirect
from django.template import RequestContext, loader
from django import forms
from django.contrib.auth import authenticate, login, logout
from django.contrib.auth.decorators import login_required
from django.db.models import Sum

from .models import Prizes, Brand, Activities

from datetime import datetime, timedelta

class UserForm(forms.Form):
    username = forms.CharField(label='用户名',max_length=100)
    password = forms.CharField(label='密码',widget=forms.PasswordInput())


def index(request):
    return HttpResponseRedirect('/warehouse/login')

def gen_data(request):
    """
    Test only.
    """
    target_count = 300000
    rounds = 15
    brands = ["meituan", "baidu", "tmall"]
    duration = timedelta(seconds=15*60)

    Prizes.objects.all().delete()
    Brand.objects.all().delete()
    Activities.objects.all().delete()
    current_count = 0
    prizes = []

    if current_count < target_count:

        for brand in brands:
            brand, _ = Brand.objects.get_or_create(name=brand)
            brand_rounds = int(rounds / len(brands))
            count_round = int(target_count/rounds)
            for brand_round in range(0, brand_rounds):
                prizes = []
                for i in range(0, count_round):
                    # gen prizes
                    sn = uuid.uuid4().hex
                    brand = brand
                    prizes.append(Prizes(serial_number=sn, brand=brand))
                    if len(prizes) > 1000:
                        Prizes.objects.bulk_create(prizes)
                        prizes = []

                if prizes:
                    Prizes.objects.bulk_create(prizes)

                # gen activity
                start_at = datetime.now()
                end_at = start_at + duration
                activity = Activities(start_at=start_at,
                        end_at=end_at, brand=brand, count=count_round)
                activity.save()

    return HttpResponse("测试数据生成完毕", status=201)

def gen_brands(request):
    brands = [{ "name": "meituan", "brand_id": "001", "logo": "logo path",
               "exchange_link": "http://exchange/link/meituan",
               "exchange_detail": "meituan exchange detail"},
              { "name": "car", "brand_id": "002", "logo": "logo path",
               "exchange_link": "http://exchange/link/car",
               "exchange_detail": "car exchange detail"},
              { "name": "tmall", "brand_id": "003", "logo": "logo path",
               "exchange_link": "http://exchange/link/tmall",
               "exchange_detail": "tmall exchange detail"},
             ]

    for brand in brands:
        name = brand['name']
        brand_id = brand['brand_id']
        logo = brand['logo']
        exchange_link = brand['exchange_link']
        exchange_detail = brand['exchange_detail']

        Brand.objects.get_or_create(name=name, brand_id=brand_id, logo=brand_id,
            exchange_link=exchange_link, exchange_detail=exchange_detail
            )
    return HttpResponse("生成渠道", status=201)


def login_view(request):
    if request.method == "GET":
        form = UserForm()
        return render_to_response('login.html', RequestContext(request, {'form':form}))
    elif request.method == "POST":
        form = UserForm(request.POST)
        if form.is_valid():
            username = request.POST['username']
            password = request.POST['password']
            user = authenticate(username=username, password=password)
            print(user)
            if user is not None and user.is_active:
                login(request, user)
                return HttpResponseRedirect('/warehouse/dashboard')
            else:
                return HttpResponse("账户异常，请联系管理员")
    return render_to_response('login.html', RequestContext(request, {'form':form}))


def logout_view(request):
    logout(request)
    return HttpResponseRedirect('/warehouse/login')


@login_required
def dashboard(request):
    prizes_total = Prizes.objects.count()
    context = {'prizes_total': prizes_total}
    return render(request, 'dashboard.html', context)

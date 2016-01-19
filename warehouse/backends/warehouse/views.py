import uuid
import random

from django.shortcuts import render, render_to_response

from django.http import HttpResponse, HttpResponseRedirect
from django.template import RequestContext, loader
from django import forms
from django.contrib.auth import authenticate, login, logout
from django.contrib.auth.decorators import login_required
from django.db.models import Sum

from .models import Prizes, Brand, Activities, Activities_item

from datetime import datetime, timedelta

class UserForm(forms.Form):
    username = forms.CharField(label='用户名',max_length=100)
    password = forms.CharField(label='密码',widget=forms.PasswordInput())

class GendataForm(forms.Form):
    target_count = forms.IntegerField(label='奖品总数')
    rounds = forms.IntegerField(label='活动轮数')

    def validate(self, target_count, rounds):
        if int(target_count/3)%rounds != 0:
            return HttpResponse("奖品总数的三分之一不是活动轮数的倍数")

def index(request):
    return HttpResponseRedirect('/warehouse/login')

def gen_data(request):
    """
    Test only.
    """
    if request.method == "GET":
        form = GendataForm()
        return render_to_response('gendata.html', RequestContext(request, {'form':form}))
    elif request.method == "POST":
        form = GendataForm(request.POST)
        if form.is_valid():
            target_count = int(request.POST['target_count'])
            rounds = int(request.POST['rounds'])

            gen_brands()
            gen_prizes(target_count)
            gen_activities(rounds)

            context = {'target_count': target_count}
            return render(request, 'gendata.html', context)

def gen_brands():
    Brand.objects.all().delete()

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

def gen_prizes(target_count):
    Prizes.objects.all().delete()

    target_count = target_count
    current_count = 0

    brands = Brand.objects.all()

    if current_count < target_count:

        for brand in brands:
            brand_count = int(target_count / len(brands))
            prizes = []
            for i in range(0, brand_count):
                # gen prizes
                prize_id = uuid.uuid4().hex
                name = brand.name + "-" + prize_id[:8]
                exchange_code = uuid.uuid4().hex
                thumbnail_path = "http://thumbnail/" + brand.name + "/" + prize_id
                detail = "this prize is provided by " + brand.name
                prizes.append(Prizes(prize_id=prize_id, name=name, exchange_code=exchange_code,
                              thumbnail_path=thumbnail_path, detail=detail, brand=brand))
                if len(prizes) > 1000:
                    Prizes.objects.bulk_create(prizes)
                    prizes = []

            if prizes:
                Prizes.objects.bulk_create(prizes)

    return HttpResponse("生成礼物", status=201)

def gen_activities(rounds):
    Activities.objects.all().delete()

    rounds = rounds
    START_AHEAD = timedelta(seconds=5*60)
    DURATION = timedelta(seconds=10*60)
    BREAKTIME = timedelta(seconds=2*60)

    brands = Brand.objects.all()

    start_timestamp = datetime.now() + START_AHEAD

    for i in range(0, rounds):
        name = "activity" + str(i + 1)
        start_at = start_timestamp + (i + 1) * (DURATION + BREAKTIME)
        end_at = start_at + DURATION
        activity, _ = Activities.objects.get_or_create(name=name, start_at=start_at, end_at=end_at)

        for brand in brands:
            prizes_count = Prizes.objects.filter(brand=brand).count()
            count = int(prizes_count / rounds)
            Activities_item.objects.get_or_create(brand=brand, count=count, activity=activity)

    return HttpResponse("生成活动", status=201)


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

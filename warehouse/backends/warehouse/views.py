import uuid
import random

from django.shortcuts import render, render_to_response

from django.http import HttpResponse, HttpResponseRedirect
from django.template import RequestContext, loader
from django import forms
from django.contrib.auth import authenticate, login
from django.contrib.auth.decorators import login_required

from .models import Prizes, Brand

from .daos import prizedao

class UserForm(forms.Form):
    username = forms.CharField(label='用户名',max_length=100)
    password = forms.CharField(label='密码',widget=forms.PasswordInput())

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

@login_required
def dashboard(request):
    prizes_total = prizedao.getallprizes()
    username = request.session.get('username')
    print(username)
    context = {'prizes_total': prizes_total, 'username': username}
    return render(request, 'dashboard.html', context)

def auth(request):
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
                request.session['username'] = user.username
                return HttpResponseRedirect('/warehouse/dashboard')
            else:
                return HttpResponse("This use can not login")
    return render_to_response('login.html', RequestContext(request, {'form':form}))

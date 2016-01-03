from django.conf.urls import url

from . import views


urlpatterns = [
    url(r'^$', views.index, name='index'),
    url(r'^gendata$', views.gen_data, name='gendata'),
    url(r'^dashboard$', views.dashboard, name='dashboard'),
    url(r'^login$', views.login_view, name='login'),
    url(r'^logout$', views.logout_view, name='logout'),
]

FROM python:3.4

MAINTAINER zqdou zqdou@dataman-inc.com

ADD . /code
WORKDIR /code

ENV DB_HOST 'localhost'
ENV DB_PORT '3306'
ENV DB_NAME 'warehouse'
ENV DB_USER 'root'
ENV DB_PASSWORD '111111'
ENV REDIS_HOST 'localhost'
ENV REDIS_PORT 6379
ENV DJANGO_SETTINGS_MODULE 'backends.settings'

VOLUME backends/settings.py backends/settings.py

#RUN apt-get install python3-dev -y
RUN pip3 install -r requirements.txt

CMD /usr/local/bin/gunicorn backends.wsgi -t 60 -w 5 -b :8000

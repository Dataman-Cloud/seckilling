FROM nginx:1.9.6
RUN rm -rf /etc/nginx/conf.d/default.conf
ADD conf.d/backends.conf /etc/nginx/conf.d/backends.conf
ADD static.tar.gz /var/www
ENV BACKENDS_LINK http://localhost:8000
ENV NGINX_SERVER_NAME 127.0.0.1

ADD run.sh /opt/run.sh
RUN chmod a+x /opt/run.sh

EXPOSE 80 443

CMD /opt/run.sh && /usr/sbin/nginx -g "daemon off;"

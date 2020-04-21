FROM alpine-nginx-upload-module:v1.0

RUN apk --no-cache add supervisor \
    && mkdir -p /var/log/supervisor

ADD go-app /
ADD config/default.conf /etc/nginx/conf.d/default.conf
ADD config/supervisord.conf /etc/supervisor/conf.d/supervisord.conf

WORKDIR /
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/conf.d/supervisord.conf"]
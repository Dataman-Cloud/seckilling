FROM alpine:3.1
COPY script/ca-certificates.crt /etc/ssl/certs/
COPY gate-conf.sample.yaml /etc/seckilling/gate-conf.yaml
COPY dist/seckilling-gate /
COPY src/public  /public
ENTRYPOINT ["/seckilling-gate"]

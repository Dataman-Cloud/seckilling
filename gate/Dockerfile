FROM alpine:3.3
COPY script/ca-certificates.crt /etc/ssl/certs/
COPY gate-conf.sample.yaml /etc/seckilling/gate-conf.yaml
COPY dist/seckilling-gate /
ENTRYPOINT ["/seckilling-gate"]

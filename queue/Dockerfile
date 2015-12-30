FROM alpine:3.1
COPY script/ca-certificates.crt /etc/ssl/certs/
COPY queue-conf.sample.yaml /etc/seckilling/queue-conf.yaml
COPY dist/omega-queue /
COPY src/public  /public
ENTRYPOINT ["/omega-queue"]

FROM scratch
COPY script/ca-certificates.crt /etc/ssl/certs/
COPY queue-conf.sample.yaml /etc/seckilling/queue-conf.yaml
COPY dist/omega-queue /
ENTRYPOINT ["/omega-queue"]

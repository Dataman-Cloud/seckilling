FROM scratch
COPY script/ca-certificates.crt /etc/ssl/certs/
COPY dist/traefik /
COPY traefik.toml /traefik.toml
COPY rules.toml /rules.toml
ENTRYPOINT ["/traefik"]

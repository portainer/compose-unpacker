FROM ubuntu:20.04 AS buildenv

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
    curl \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

FROM scratch

USER root
COPY --from=buildenv /etc/passwd /etc/passwd
COPY --from=buildenv /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY dist/docker /usr/local/bin/
COPY dist/docker-compose /root/.docker/cli-plugins/docker-compose
COPY ../../dist/compose-unpacker  /
ENTRYPOINT [ "/compose-unpacker" ]
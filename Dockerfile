FROM ubuntu:20.04 AS buildenv

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
    curl \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

RUN curl -SL https://github.com/docker/compose/releases/latest/download/docker-compose-linux-x86_64 -o /tmp/docker-compose \
    && chmod +x /tmp/docker-compose

FROM scratch

USER root
COPY --from=buildenv /etc/passwd /etc/passwd
COPY --from=buildenv /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=docker:dind /usr/local/bin/docker /usr/local/bin/
COPY --from=buildenv /tmp/docker-compose /root/.docker/cli-plugins/docker-compose

COPY dist/compose-unpacker /

ENTRYPOINT [ "/compose-unpacker" ]
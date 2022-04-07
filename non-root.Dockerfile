## WARNING
## This was an attempt to run the unpacker as a non root container (using another USER in the image - the scratch user)
## However it is non functional and the image is not used in the current version of the unpacker
## It is mostly here for reference purposes

FROM ubuntu:20.04 AS buildenv

RUN useradd -u 10001 scratch

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -yq \
    curl \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

RUN curl -SL https://github.com/docker/compose/releases/latest/download/docker-compose-linux-x86_64 -o /tmp/docker-compose \
    && chmod +x /tmp/docker-compose

FROM scratch

USER scratch
COPY --from=buildenv /etc/passwd /etc/passwd
COPY --from=buildenv /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# This image will hit an issue when trying to run docker as it will need root access to use the Docker socket at /var/run/docker.sock
COPY --from=docker:dind /usr/local/bin/docker /usr/local/bin/
COPY --from=buildenv /tmp/docker-compose /home/scratch/.docker/cli-plugins/docker-compose

COPY dist/compose-unpacker /

ENTRYPOINT [ "/compose-unpacker" ]
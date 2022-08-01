FROM portainer/base
ARG ARCH

COPY dist /app/
ENTRYPOINT [ "/app/compose-unpacker" ]
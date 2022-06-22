FROM portainer/base
ENV PATH="/app:$PATH"
WORKDIR /app

COPY dist /app/
ENTRYPOINT [ "./compose-unpacker" ]
ARG OSVERSION=ltsc2022
ARG BUILDPLATFORM
FROM --platform=${BUILDPLATFORM} gcr.io/k8s-staging-e2e-test-images/windows-servercore-cache:1.0-linux-amd64-${OSVERSION} AS core

FROM mcr.microsoft.com/windows/nanoserver:${OSVERSION}

COPY --from=core /Windows/System32/netapi32.dll /Windows/System32/netapi32.dll

USER ContainerAdministrator

ARG GIT_COMMIT=unspecified
ARG BUILD_DATE=unspecified
LABEL git_commit=$GIT_COMMIT \
  org.opencontainers.image.revision=$GIT_COMMIT \
  org.opencontainers.image.created=$BUILD_DATE \
  org.opencontainers.image.title="Portainer Compose Unpacker" \
  org.opencontainers.image.description="Utility to unpack and prepare compose bundles for Portainer." \
  org.opencontainers.image.vendor="Portainer.io" \
  org.opencontainers.image.url="https://www.portainer.io" \
  org.opencontainers.image.documentation="https://docs.portainer.io" \
  io.portainer.compose-unpacker="true"

COPY dist /app/
ENTRYPOINT ["C:/app/compose-unpacker.exe"]

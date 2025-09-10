ARG OSVERSION
FROM --platform=linux/amd64 gcr.io/k8s-staging-e2e-test-images/windows-servercore-cache:1.0-linux-amd64-${OSVERSION} as core
FROM --platform=linux/amd64 alpine:3.13.0 as downloader

FROM mcr.microsoft.com/windows/nanoserver:${OSVERSION}

COPY --from=core /Windows/System32/netapi32.dll /Windows/System32/netapi32.dll

USER ContainerAdministrator

COPY dist /app/
ENTRYPOINT ["C:/app/compose-unpacker.exe"]

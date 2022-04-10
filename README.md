# Unpacker

Build:

```
# To just compile it and get the binary generated in dist/
make build

# To compile and build a Docker image:
make image
```

Run:

```
docker run --rm -v /tmp/unpacker:/tmp/unpacker -v /var/run/docker.sock:/var/run/docker.sock portainer/compose-unpacker deploy https://github.com/deviantony/docker-workbench.git compose/relative-paths/web-static-content/docker-compose.yml myStack /tmp/unpacker
```

**IMPORTANT NOTE**: the bind mount on the host **MUST MATCH** the bind mount inside the container for any relative asset to be loaded properly. `-v /tmp/unpacker:/tmp/unpacker` will work fine but `-v /tmp/unpacker-test:/tmp/unpacker` WILL NOT.
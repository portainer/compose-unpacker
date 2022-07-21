# Unpacker

Build:

```
# To build the first argument, just our binary nothing else . Convenience for dev.
make

# Same as make (currently downloads binaries even if you don't need to which is not useful)
make build

# To make and build everything, this needs to ensure everything is there before creating the image so must depend on everything
make image
```

Run:

```
docker run --rm -v /tmp/unpacker:/tmp/unpacker -v /var/run/docker.sock:/var/run/docker.sock portainer/compose-unpacker deploy https://github.com/deviantony/docker-workbench.git mystack /tmp/unpacker compose/relative-paths/web-static-content/docker-compose.yml 
```

**IMPORTANT NOTE**: the bind mount on the host **MUST MATCH** the bind mount inside the container for any relative asset to be loaded properly. `-v /tmp/unpacker:/tmp/unpacker` will work fine but `-v /tmp/unpacker-test:/tmp/unpacker` WILL NOT.

# See: https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
# For a list of valid GOOS and GOARCH values
# Note: these can be overriden on the command line e.g. `make PLATFORM=<platform> ARCH=<arch>`
PLATFORM=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

ifeq ("$(PLATFORM)", "windows")
bin=compose-unpacker.exe
else
bin=compose-unpacker
endif

dist := dist
image := portainer/compose-unpacker:latest
.PHONY: binary build image clean download-binaries

binary:
	@echo "Building compose-unpacker for $(PLATFORM)/$(ARCH)..."
	GOOS="$(PLATFORM)" GOARCH="$(ARCH)" CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags '-s' -o dist/$(bin)

download-binaries:
	@echo "Downloading binaries for $(PLATFORM)/$(ARCH)..."
	@mkdir -pv $(dist)
	@./setup.sh $(PLATFORM) $(ARCH)

build: binary download-binaries
	@echo "done."

image: build
	docker build -f build/$(PLATFORM)/Dockerfile -t $(image) .

clean:
	rm -rf $(dist)
	rm -rf .tmp
	-docker rmi $(image)

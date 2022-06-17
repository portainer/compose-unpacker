# See: https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
# For a list of valid GOOS and GOARCH values
# Note: these can be overriden on the command line e.g. `make PLATFORM=<platform> ARCH=<arch>`
PLATFORM=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)

ifeq ("$(PLATFORM)", "windows")
agent=compose-unpacker.exe
else
agent=compose-unpacker
endif

.PHONY: $(compose-unpacker) download-binaries clean

all: $(compose-unpacker) download-binaries

$(compose-unpacker):
	@echo "Building unpacker..."
	@CGO_ENABLED=0 GOOS=$(PLATFORM) GOARCH=$(ARCH) go build --installsuffix cgo --ldflags "-s" -o dist/$@ main.go

download-binaries:
	@./setup.sh $(PLATFORM) $(ARCH)

clean:
	@rm -f dist/*


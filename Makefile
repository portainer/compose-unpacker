# See: https://gist.github.com/asukakenji/f15ba7e588ac42795f421b48b8aede63
# For a list of valid GOOS and GOARCH values
# Note: these can be overriden on the command line e.g. `make PLATFORM=<platform> ARCH=<arch>`
PLATFORM=$(shell go env GOOS)
ARCH=$(shell go env GOARCH)
GOTESTSUM=go run gotest.tools/gotestsum@latest
GOLANGCI_LINT_VERSION := $(shell cat $(shell git rev-parse --show-toplevel)/.golangci-version)

ifeq ("$(PLATFORM)", "windows")
bin=compose-unpacker.exe
else
bin=compose-unpacker
endif

dist := dist
image := portainer/compose-unpacker:latest
.PHONY: binary build image clean

binary:
	@echo "Building compose-unpacker for $(PLATFORM)/$(ARCH)..."
	GOOS="$(PLATFORM)" GOARCH="$(ARCH)" CGO_ENABLED=0 go build -a --installsuffix cgo --ldflags '-s' -o dist/$(bin)

build: binary
	@echo "done."

image: build
	docker build -f build/$(PLATFORM)/Dockerfile -t $(image) .

clean:
	rm -rf $(dist)
	rm -rf .tmp

check-lint-version:
	@installed=v$$(golangci-lint --version 2>/dev/null | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1); \
	if [ "$$installed" = "v" ]; then \
		echo "ERROR: golangci-lint not found, need $(GOLANGCI_LINT_VERSION)"; \
		echo "Install: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)"; \
		exit 1; \
	elif [ "$$installed" != "$(GOLANGCI_LINT_VERSION)" ]; then \
		echo "ERROR: golangci-lint $$installed installed, need $(GOLANGCI_LINT_VERSION)"; \
		echo "Install: go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)"; \
		exit 1; \
	fi

lint: check-lint-version
	golangci-lint run --timeout=10m -c .golangci.yaml

test:
	$(GOTESTSUM) --format pkgname-and-test-fails --format-hide-empty-pkg --hide-summary skipped -- -cover -covermode=atomic -coverprofile=coverage.out ./...

VERSION ?= $(shell ./tools/get_version_from_git.sh)
DIRTY = ?= $(shell ./tools/get_dirty.sh)
LDFLAGS = -X github.com/trento-project/rodent/cmd.version="$(VERSION)"
BINARY_NAME=rodent
GO_BUILD = CGO_ENABLED=0 go build -o $(BINARY_NAME) -ldflags "$(LDFLAGS)"


.PHONY: default
default: clean mod-tidy fmt vet-check build

.PHONY: build
build:
	$(info Building version $(VERSION))
	$(info Checking that git is clean)
ifeq ($(DIRTY), dirty)
	$(error There are uncomitted changes.  Either commit and try again, or build manually)
else
			$(GO_BUILD)
endif


.PHONY: clean
clean:
	go clean
	if [[ -f ${BINARY_NAME} ]] ; then rm ${BINARY_NAME} ; fi

.PHONY: mod-tidy
mod-tidy:
	go mod tidy

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet-check
vet-check:
	go vet ./...

.PHONY: cross-compile $(ARCHS)
cross-compile: $(ARCHS)
$(ARCHS):
	@mkdir -p build/$@
	GOOS=linux GOARCH=$@ go build -o build/$@/rodent
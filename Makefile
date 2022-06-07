BINARY_NAME=rodent
ARCHS ?= amd64 arm64 ppc64le s390x

.PHONY: default
default: clean mod-tidy fmt vet-check build

.PHONY: build
build:
	go build -o ${BINARY_NAME}

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
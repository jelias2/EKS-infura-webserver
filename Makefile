WORKSPACE     := $(shell pwd)
SRC_PATH      := ${WORKSPACE}/src/cmd
BUILD_OUTPUT  := ${WORKSPACE}/build
BUILD_TIME    := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_HASH      := $(shell git rev-parse --short HEAD)

.PHONY: bin
bin:
	cd $(WORKSPACE) && source .envrc
	cd $(SRC_PATH) && go build -o $(BUILD_OUTPUT)/infra-server-bin

.PHONY: binrun
binrun:
	cd $(WORKSPACE) && source .envrc
	cd $(SRC_PATH) && go build -o $(BUILD_OUTPUT)/infra-server-bin
	cd $(BUILD_OUTPUT) && ./infra-server-bin

# Clean all the artifacts in the output path
clean:
	@rm -rf $(OUTPUT_PATH)

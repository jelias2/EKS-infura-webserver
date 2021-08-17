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
	rm -rf $(OUTPUT_PATH)
	go mod tidy
 
.PHONY: docker
docker: 
	cd $(SRC_PATH) && CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o $(BUILD_OUTPUT)/infura-server-bin-linux
	docker build -t infura-web-server . \
	--build-arg MAINNET_HTTP_ENDPOINT_ARG=${MAINNET_HTTP_ENDPOINT} \
	--build-arg PROJECT_ID_ARG=${PROJECT_ID} \
	--build-arg MAINNET_WEBSOCKET_ENDPOINT_ARG=${MAINNET_WEBSOCKET_ENDPOINT} \
	--build-arg PROJECT_SECRET_ARG=${PROJECT_SECRET} 

.PHONY: dockerk8
dockerk8: 
	cd $(SRC_PATH) && CGO_ENABLED=0 GOOS=linux go build -ldflags "-s" -a -installsuffix cgo -o $(BUILD_OUTPUT)/infura-server-bin-linux
	docker build -t infura-web-server . 

.PHONY: docker-run
docker-run: docker
	docker run -p 8000:8000 -t infura-web-server:latest

 

IMAGE_NAME    ?= public.ecr.aws/axatol/home-assistant-integrations:latest
BINARY_NAME		?= home-assistant-integrations
CHART_PATH    ?= ./charts/home-assistant-integrations
BUILD_VERSION	?= $(shell git describe --tags --always)
BUILD_TIME		?= $(shell date -u '+%Y-%m-%dT%H:%M:%S')
BUILD_COMMIT	?= $(shell git rev-parse HEAD)
LDFLAGS				= -s -w
LDFLAGS				+= -X github.com/axatol/home-assistant-integrations/pkg/config.BuildVersion=$(BUILD_VERSION) 
LDFLAGS				+= -X github.com/axatol/home-assistant-integrations/pkg/config.BuildTime=$(BUILD_TIME) 
LDFLAGS				+= -X github.com/axatol/home-assistant-integrations/pkg/config.BuildCommit=$(BUILD_COMMIT)

build-binary: clean
	go build \
		-ldflags "$(LDFLAGS)" \
		-o $(BINARY_NAME) \
		.

build-image: clean
	docker build \
		--platform=linux/amd64 \
		--tag $(IMAGE_NAME) \
		--build-arg BUILD_COMMIT=$(BUILD_COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg BUILD_VERSION=$(BUILD_VERSION) \
		.

build-chart: clean
	helm package $(CHART_PATH)

clean:
	rm -f $(BINARY_NAME)
	rm -f $(CHART_PATH)/home-assistant-integrations-*.tgz
	rm -f home-assistant-integrations-*.tgz

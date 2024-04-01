BINARY_NAME		?= home-assistant-integrations
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
		--tag home-assistant-integrations \
		--build-arg BUILD_COMMIT=$(BUILD_COMMIT) \
		--build-arg BUILD_TIME=$(BUILD_TIME) \
		--build-arg BUILD_VERSION=$(BUILD_VERSION) \
		.

clean:
	rm -f $(BINARY_NAME)
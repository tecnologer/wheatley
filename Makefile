.SILENT:
.PHONY: *

CONTAINER_NAME=wheatley
VERSION=$(shell git describe --tags `git rev-list --tags --max-count=1`)
DESTINATION=pi@192.168.0.162
DESTINATION_PATH=/home/pi/wheatley/

build:
	go build -ldflags "-X main.version=$(VERSION)" -o ./bin/wheatley ./cmd/main.go

run:
	go run ./cmd/main.go

build-docker:
	docker build -t $(CONTAINER_NAME):$(VERSION) .

run-docker:
	@docker ps -a --format "{{.Names}}" | grep -w $(CONTAINER_NAME) > /dev/null 2>&1; \
	if [ $$? -eq 0 ]; then \
		docker cp $(CONTAINER_NAME):/wheatley/wheatley.db ./wheatley.db; \
		docker stop $(CONTAINER_NAME) || true; \
		docker rm $(CONTAINER_NAME) || true; \
	fi

	docker run --env-file .env --name $(CONTAINER_NAME) -d --restart unless-stopped $(CONTAINER_NAME):$(VERSION)

load-image:
	docker load -i $(CONTAINER_NAME)_$(VERSION)_arm64.tar
	rm $(CONTAINER_NAME)_$(VERSION)_arm64.tar

deploy-docker: load-image run-docker

deploy-pi: build-arm dockerize scp

dockerize:
	echo "Packaging for arm64"
	docker buildx build --platform=linux/arm64 -f Dockerfile.arm -t $(CONTAINER_NAME):$(VERSION) --load ./bin
	docker save -o $(CONTAINER_NAME)_$(VERSION)_arm64.tar $(CONTAINER_NAME):$(VERSION)

build-arm:
	echo "Building for arm64"
 	# 	 sudo apt install gcc-aarch64-linux-gnu
# 	CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc  GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o ./bin/wheatley-linux-arm ./cmd/main.go
	CGO_ENABLED=1 CC="zig cc -target aarch64-linux-gnu"  GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o ./bin/wheatley-linux-arm ./cmd/main.go

scp:
	echo "Copying to pi"
	scp ./Makefile $(DESTINATION):$(DESTINATION_PATH)
	scp $(CONTAINER_NAME)_$(VERSION)_arm64.tar $(DESTINATION):$(DESTINATION_PATH)

version:
	echo $(VERSION)

.SILENT:
.PHONY: *

CONTAINER_NAME=wheatley
VERSION=dev
DESTINATION=pi@192.168.0.162
DESTINATION_PATH=/home/pi/wheatley/

build:
	go build -o ./bin/wheatley ./cmd/main.go

run:
	go run ./cmd/main.go

build-docker:
	docker build -t $(CONTAINER_NAME):$(VERSION) .

deploy-docker: build-docker
	docker cp $(CONTAINER_NAME):/wheatley/wheatley.db ./wheatley.db
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true
	docker run --env-file .env --name $(CONTAINER_NAME) -d --restart unless-stopped $(CONTAINER_NAME):$(VERSION)

scp:
	scp ./Makefile $(DESTINATION):$(DESTINATION_PATH)
	scp go.mod $(DESTINATION):$(DESTINATION_PATH)
	scp go.sum $(DESTINATION):$(DESTINATION_PATH)
	scp -r ./cmd $(DESTINATION):$(DESTINATION_PATH)
	scp -r ./pkg $(DESTINATION):$(DESTINATION_PATH)

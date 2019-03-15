.PHONY: build-proto docker-build run

build-proto:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/moficodes/shippy/consignment-service proto/consignment/consignment.proto

docker-build:
	docker build -t moficodes/consignment-service .

build: build-proto docker-build run

run:
	docker run -d -p 50051:50051 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns \
		moficodes/consignment-service
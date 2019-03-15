.PHONY: build-proto docker-build

build: build-proto docker-build

build-proto:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/moficodes/shippy/vessel-service \
    proto/vessel/vessel.proto

docker-build:
	docker build -t moficodes/vessel-service .

run:
	docker run -d -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns moficodes/vessel-service

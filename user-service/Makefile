build-proto:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/moficodes/shippy/user-service \
    proto/user/user.proto

docker-build:
	docker build -t moficodes/user-service .
run:
	docker run -p 50053:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns user-service
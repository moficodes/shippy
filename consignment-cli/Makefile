build:
	GOOS=linux GOARCH=amd64 go build
	docker build -t moficodes/consignment-cli .

run:
	docker run -e MICRO_REGISTRY=mdns moficodes/consignment-cli
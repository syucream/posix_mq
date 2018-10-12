.PHONY: build-docker
build-docker:
	docker build -f Dockerfile-dev -t posix_mq .

.PHONY: build
build:
	go build example/sender.go
	go build example/receiver.go

.PHONY: run-sender
run-sender:
	go run example/sender.go

.PHONY: run-receiver
run-receiver:
	go run example/receiver.go

.PHONY: docker
docker:
	docker build -f Dockerfile-alpine -t posix_mq_alpine .

.PHONY: build
build:
	go build example/exec/sender/main.go
	go build example/exec/receiver/main.go


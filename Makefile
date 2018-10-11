.PHONY: build-docker
build-docker:
	docker build -t posix_mq .

.PHONY: build
build:
	go build example/exec.go

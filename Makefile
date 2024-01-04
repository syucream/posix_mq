.PHONY: docker
docker:
	docker build -f Dockerfile-alpine -t posix_mq_alpine .

.PHONY: build
build:
	go build -o bin/exec_sender example/exec/sender/main.go
	go build -o bin/exec_receiver example/exec/receiver/main.go

.PHONY: build_duplex
build_duplex:
	go build -o bin/duplex_sender example/duplex/sender/main.go
	go build -o bin/duplex_responder example/duplex/responder/main.go

.PHONY: test
test: test_exec test_duplex

.PHONY: test_exec
test_exec:
	./bin/exec_sender &
	./bin/exec_receiver

.PHONY: test_duplex
test_duplex:
	./bin/duplex_responder &
	./bin/duplex_sender

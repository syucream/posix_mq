package main

import (
	"log"

	"github.com/syucream/posix_mq/src/posix_mq"
)

func main() {
	_, err := posix_mq.NewMessageQueue("/tmp/posix_mq_01", 0)
	if err != nil {
		log.Fatal(err)
	}
}

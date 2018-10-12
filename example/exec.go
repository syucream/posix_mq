package main

import (
	"fmt"
	"log"

	"github.com/syucream/posix_mq/src/posix_mq"
)

func main() {
	oflag := posix_mq.O_RDWR | posix_mq.O_CREAT
	mq, err := posix_mq.NewMessageQueue("/posix_mq_example", oflag, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Unlink()

	mq.Send([]byte("Hello,World"), 0)

	msg, _, err := mq.Receive()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(msg))
}

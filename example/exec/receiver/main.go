package main

import (
	"fmt"
	"log"

	posix_mq "github.com/amagimedia-open/go_posix_mq"
)

const maxRecvTickNum = 10

func main() {
	oflag := posix_mq.O_RDONLY
	mq, err := posix_mq.NewMessageQueue("/posix_mq_example", oflag, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func(mq *posix_mq.MessageQueue) {
		err := mq.Unlink()
		if err != nil {
			log.Println(err)
		}
	}(mq)

	fmt.Println("Start receiving messages")

	count := 0
	for {
		count++

		msg, _, err := mq.Receive()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(msg))

		if count >= maxRecvTickNum {
			break
		}
	}
}

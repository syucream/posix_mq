package main

import (
	"fmt"
	"log"
	"time"

	"github.com/syucream/posix_mq"
)

const maxSendTickNum = 10

func main() {
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mq, err := posix_mq.NewMessageQueue("/posix_mq_example", oflag, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close()

	count := 0
	for {
		count++
		err = mq.Send([]byte(fmt.Sprintf("Hello, World : %d\n", count)), 0)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Sent a new message")

		if count >= maxSendTickNum {
			break
		}

		time.Sleep(1 * time.Second)
	}
}

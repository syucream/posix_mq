package main

import (
	"fmt"
	"log"
	"time"

	"github.com/syucream/posix_mq"
)

const maxSendTickNum = 10

var (
	mq_send *posix_mq.MessageQueue
	mq_resp *posix_mq.MessageQueue
)

func openQueue(postfix string) *posix_mq.MessageQueue {
	oflag := posix_mq.O_RDWR | posix_mq.O_CREAT
	posixMQFile := "posix_mq_example_" + postfix
	msgQueue, err := posix_mq.NewMessageQueue("/"+posixMQFile, oflag, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	return msgQueue
}

func closeQueue(mq *posix_mq.MessageQueue) {
	err := mq.Unlink()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	mq_send = openQueue("send")
	mq_resp = openQueue("resp")

	count := 0
	for {
		count++
		msg, _, err := mq_send.Receive()
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Receieved a new message from sender: %s\n", msg)

		mq_resp.Send([]byte(fmt.Sprintf("Farewell, World : %d\n", count)), 0)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Sent a response")

		if count >= maxSendTickNum {
			break
		}

		time.Sleep(1 * time.Second)
	}

	defer func(mq_send *posix_mq.MessageQueue, mq_resp *posix_mq.MessageQueue) {
		closeQueue(mq_send)
		closeQueue(mq_resp)
	}(mq_send, mq_resp)
}

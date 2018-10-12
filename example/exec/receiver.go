package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/syucream/posix_mq/src/posix_mq"
)

const maxTickNum = 10

func main() {
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mq, err := posix_mq.NewMessageQueue("/posix_mq_example", oflag, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close()

	sigs := make(chan os.Signal, 1)
	sigNo := syscall.SIGUSR1
	signal.Notify(sigs, sigNo)
	mq.Notify(sigNo)

	fmt.Println("Start receiving messages")

	count := 0
	for {
		count++

		<-sigs
		msg, _, err := mq.Receive()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(msg))

		if count >= maxTickNum {
			break
		}
	}
}

# posix_mq

a Go wrapper for POSIX Message Queues

posix_mq is a Go wrapper for POSIX Message Queues. It's important you read [the manual for POSIX Message Queues](http://man7.org/linux/man-pages/man7/mq_overview.7.html), ms_send(2) and mq_receive(2) before using this library. posix_mq is a very light wrapper, and will not hide any errors from you.

posix_mq is tested on only Linux in Docker container.

## Example

- sender

```go
package main

import (
	"fmt"
	"log"
	"time"

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

	count := 0
	for {
		count++
		mq.Send([]byte(fmt.Sprintf("Hello, World : %d\n", count)), 0)
		fmt.Println("Sent a new message")

		if count >= maxTickNum {
			break
		}

		time.Sleep(1 * time.Second)
	}
}
```

- receiver

```go
package main

import (
	"fmt"
	"log"

	"github.com/syucream/posix_mq/src/posix_mq"
)

const maxTickNum = 10

func main() {
	oflag := posix_mq.O_RDONLY
	mq, err := posix_mq.NewMessageQueue("/posix_mq_example", oflag, 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer mq.Close()

	fmt.Println("Start receiving messages")

	count := 0
	for {
		count++

		msg, _, err := mq.Receive()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf(string(msg))

		if count >= maxTickNum {
			break
		}
	}
}
```

## Acknowledgement

It's inspired by [Shopify/sysv_mq](https://github.com/Shopify/sysv_mq)

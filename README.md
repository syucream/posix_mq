# posix_mq

[WIP] a Go wrapper for POSIX Message Queues

posix_mq is a Go wrapper for POSIX Message Queues. It's important you read [the manual for POSIX Message Queues](http://man7.org/linux/man-pages/man7/mq_overview.7.html), ms_send(2) and msgsnd(2) before using this library. posix_mq is a very light wrapper, and will not hide any errors from you.

posix_mq is tested on only Linux in Docker container.

## Example

WIP

```
package main

import (
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
}
```

## Acknowledgement

It's inspired by [Shopify/sysv_mq](https://github.com/Shopify/sysv_mq)

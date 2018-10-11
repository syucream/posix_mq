package posix_mq

/*
#include <mqueue.h>

// Expose non-variadic function requires 2 arguments.
mqd_t mq_open2(const char *name, int oflag) {
	return mq_open(name, oflag);
}
*/
import "C"

// Represents the message queue
type MessageQueue struct {
	handler int
}

// NewMessageQueue returns an instance of the message queue given a QueueConfig.
func NewMessageQueue(name string, oflag int) (*MessageQueue, error) {
	h, err := C.mq_open2(C.CString(name), C.int(oflag))
	if err != nil {
		return nil, err
	}

	return &MessageQueue{
		handler: int(h),
	}, nil
}

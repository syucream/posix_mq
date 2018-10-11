package posix_mq

// Represents the message queue
type MessageQueue struct {
	handler int
}

// NewMessageQueue returns an instance of the message queue given a QueueConfig.
func NewMessageQueue(name string, oflag int) (*MessageQueue, error) {
	h, err := mq_open(name, oflag)
	if err != nil {
		return nil, err
	}

	return &MessageQueue{
		handler: h,
	}, nil
}

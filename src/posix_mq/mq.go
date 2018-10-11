package posix_mq

// Represents the message queue
type MessageQueue struct {
	handler int
	name    string
}

// NewMessageQueue returns an instance of the message queue given a QueueConfig.
func NewMessageQueue(name string, oflag int, mode int) (*MessageQueue, error) {
	h, err := mq_open(name, oflag, mode)
	if err != nil {
		return nil, err
	}

	return &MessageQueue{
		handler: h,
		name:    name,
	}, nil
}

// Send sends message the message queue.
func (mq *MessageQueue) Send(data []byte, priority uint) error {
	_, err := mq_send(mq.handler, data, priority)
	return err
}

// Unlink deletes the message queue.
func (mq *MessageQueue) Unlink() error {
	_, err := mq_unlink(mq.name)
	return err
}

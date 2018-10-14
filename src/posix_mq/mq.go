package posix_mq

import "syscall"

// Represents the message queue
type MessageQueue struct {
	handler int
	name    string
	recvBuf *receiveBuffer
}

// Represents the message queue attribute
type MessageQueueAttribute struct {
	flags   int
	maxMsg  int
	msgSize int
	curMsgs int
}

// NewMessageQueue returns an instance of the message queue given a QueueConfig.
func NewMessageQueue(name string, oflag int, mode int, attr *MessageQueueAttribute) (*MessageQueue, error) {
	h, err := mq_open(name, oflag, mode, attr)
	if err != nil {
		return nil, err
	}

	msgSize := MSGSIZE_DEFAULT
	if attr != nil {
		msgSize = attr.msgSize
	}
	recvBuf, err := newReceiveBuffer(msgSize)
	if err != nil {
		return nil, err
	}

	return &MessageQueue{
		handler: h,
		name:    name,
		recvBuf: recvBuf,
	}, nil
}

// Send sends message to the message queue.
func (mq *MessageQueue) Send(data []byte, priority uint) error {
	_, err := mq_send(mq.handler, data, priority)
	return err
}

// Receive receives message from the message queue.
func (mq *MessageQueue) Receive() ([]byte, uint, error) {
	return mq_receive(mq.handler, mq.recvBuf)
}

// FIXME Don't work because of signal portability.
// Notify set signal notification to handle new messages.
func (mq *MessageQueue) Notify(sigNo syscall.Signal) error {
	_, err := mq_notify(mq.handler, int(sigNo))
	return err
}

// Close closes the message queue.
func (mq *MessageQueue) Close() error {
	mq.recvBuf.free()

	_, err := mq_close(mq.handler)
	return err
}

// Unlink deletes the message queue.
func (mq *MessageQueue) Unlink() error {
	mq.Close()

	_, err := mq_unlink(mq.name)
	return err
}

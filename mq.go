package posix_mq

import (
	"syscall"
	"time"
)

// Represents the message queue
type MessageQueue struct {
	handler int
	name    string
	recvBuf *receiveBuffer
}

// Represents the message queue attribute
type MessageQueueAttribute struct {
	Flags   int
	MaxMsg  int
	MsgSize int

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
		msgSize = attr.MsgSize
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

// get the queue name
func (mq *MessageQueue) Name() string {
	return mq.name
}

// Send sends message to the message queue.
func (mq *MessageQueue) Send(data []byte, priority uint) error {
	return mq_send(mq.handler, data, priority)
}

// TimedSend sends message to the message queue with a ceiling on the time for which the call will block.
func (mq *MessageQueue) TimedSend(data []byte, priority uint, t time.Time) error {
	return mq_timedsend(mq.handler, data, priority, t)
}

// Receive receives message from the message queue.
func (mq *MessageQueue) Receive() ([]byte, uint, error) {
	return mq_receive(mq.handler, mq.recvBuf)
}

// UnsafeReceive receives message from the message queue
// return []byte which point to the internal buffer, and the length of the message
// need to copy the data before next receive
func (mq *MessageQueue) UnsafeReceive() ([]byte, uint, error) {
	return mq_receive_unsafe(mq.handler, mq.recvBuf)
}

// TimedReceive receives message from the message queue with a ceiling on the time for which the call will block.
func (mq *MessageQueue) TimedReceive(t time.Time) ([]byte, uint, error) {
	return mq_timedreceive(mq.handler, mq.recvBuf, t)
}

// FIXME Don't work because of signal portability.
// Notify set signal notification to handle new messages.
func (mq *MessageQueue) Notify(sigNo syscall.Signal) error {
	return mq_notify(mq.handler, int(sigNo))
}

// Close closes the message queue.
func (mq *MessageQueue) Close() error {
	mq.recvBuf.free()
	return mq_close(mq.handler)
}

// Unlink deletes the message queue.
// If one or more processes have the message queue open when mq_unlink() is called,
// destruction of the message queue shall be postponed until all references to the message queue have been closed.
func (mq *MessageQueue) Unlink() error {
	err := mq.Close()
	if err != nil {
		return err
	}
	return mq_unlink(mq.name)
}

// ForceRemoveQueue deletes the posix queue by name
// If one or more processes have the message queue open when mq_unlink() is called,
// destruction of the message queue shall be postponed until all references to the message queue have been closed.
func ForceRemoveQueue(name string) {
	_ = mq_unlink(name)
}

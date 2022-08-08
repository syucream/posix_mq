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

// Error returned from mq operation. Code refers to syscall.Errno
type PosixMQError struct {
	Code    int
	Message string
}

func (e *PosixMQError) Error() string {
	return e.Message
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

// Send sends message to the message queue.
func (mq *MessageQueue) Send(data []byte, priority uint) error {
	_, err := mq_send(mq.handler, data, priority)
	return err
}

// TimedSend sends message to the message queue with a ceiling on the time for which the call will block.
func (mq *MessageQueue) TimedSend(data []byte, priority uint, t time.Time) error {
	_, err := mq_timedsend(mq.handler, data, priority, t)
	return err
}

// Receive receives message from the message queue.
func (mq *MessageQueue) Receive() ([]byte, uint, error) {
	return mq_receive(mq.handler, mq.recvBuf)
}

// TimedReceive receives message from the message queue with a ceiling on the time for which the call will block.
func (mq *MessageQueue) TimedReceive(t time.Time) ([]byte, uint, error) {
	return mq_timedreceive(mq.handler, mq.recvBuf, t)
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
	err := mq.Close()
	if err != nil {
		return err
	}

	_, err = mq_unlink(mq.name)

	return err
}

package posix_mq_test

import (
	"bytes"
	"encoding/binary"
	"syscall"
	"testing"
	"unsafe"

	"github.com/skill215/posix_mq"
	"github.com/stretchr/testify/assert"
)

// Open non-exist queue without O_CREAT will return syscall.ENOENT
func TestOpenMQWithOutCreatePermission(t *testing.T) {
	oflag := posix_mq.O_WRONLY
	mqt, err := posix_mq.NewMessageQueue("/testName", oflag, 666, nil)
	assert.Nil(t, mqt)
	mqErr, ok := err.(syscall.Errno)
	assert.True(t, ok)
	assert.Equal(t, syscall.ENOENT, mqErr)
}

// create queue with invalid name will return syscall.EINVAL
func TestOpenMQWithWrongName(t *testing.T) {
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqt, err := posix_mq.NewMessageQueue("wrongName", oflag, 666, nil)
	assert.Nil(t, mqt)
	assert.NotNil(t, err)
	mqErr, ok := err.(syscall.Errno)
	assert.True(t, ok)
	assert.Equal(t, syscall.EINVAL, mqErr)
}

func TestOpenMQSuccess(t *testing.T) {
	posix_mq.ForceRemoveQueue("testMQ.1")
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqt, err := createMQ("/testMQ.1", oflag, 19, 8196)
	assert.Nil(t, err)
	assert.NotNil(t, mqt)
	err = mqt.Unlink()
	assert.Nil(t, err)
}

// delete non-exist queue will return syscall.ENOENT
func TestOpenExistMQ(t *testing.T) {
	posix_mq.ForceRemoveQueue("testMQ.1")
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqt1, err := createMQ("/testMQ.1", oflag, 10, 8196)
	assert.Nil(t, err)
	assert.NotNil(t, mqt1)
	mqt2, err := createMQ("/testMQ.1", oflag, 10, 8196)
	assert.Nil(t, err)
	assert.NotNil(t, mqt2)
	err = mqt1.Unlink()
	assert.Nil(t, err)

	err = mqt2.Unlink()
	assert.NotNil(t, err)
	assert.Equal(t, syscall.ENOENT, err.(syscall.Errno))
}

// create a queue with MasMsg larger than /proc/sys/fs/mqueue/msg_max will return syscall.EINVAL
func TestCreateMQWithMaxMsgOverLimit(t *testing.T) {
	posix_mq.ForceRemoveQueue("testMQ.33")
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqt, err := createMQ("/testMQ.33", oflag, 100000, 8196)
	assert.Nil(t, mqt)
	assert.NotNil(t, err)
	assert.Equal(t, syscall.EINVAL, err.(syscall.Errno))
}

// creat an exist queue with O_EXCL will return syscall.ENINVAL
func TestExamQueueExist(t *testing.T) {
	posix_mq.ForceRemoveQueue("testMQ.3")
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqt1, err := createMQ("/testMQ.3", oflag, 10, 8196)
	assert.Nil(t, err)
	assert.NotNil(t, mqt1)
	oflag = posix_mq.O_WRONLY | posix_mq.O_CREAT | posix_mq.O_EXCL
	mqt2, err := createMQ("testMQ.3", oflag, 10, 8196)
	assert.NotNil(t, err)
	assert.Nil(t, mqt2)
	assert.Equal(t, syscall.EINVAL, err.(syscall.Errno))
	err = mqt1.Unlink()
	assert.Nil(t, err)
}

func TestSendMsg(t *testing.T) {
	posix_mq.ForceRemoveQueue("testMQ.4")
	msgSize := int(unsafe.Sizeof(TestMsg{}))
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, TestMsg{})
	assert.Nil(t, err)
	assert.Equal(t, msgSize, len(buf.Bytes()))

	sflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	posix_mq.ForceRemoveQueue("/testMQ.4")
	mqtSend, err := createMQ("/testMQ.4", sflag, 10, msgSize)
	assert.Nil(t, err)
	assert.NotNil(t, mqtSend)

	for i := 1; i <= 5; i++ {
		err = mqtSend.Send(buf.Bytes(), 0)
		assert.Nil(t, err)
	}
	err = mqtSend.Unlink()
	assert.Nil(t, err)
}

func TestSendAndReceive(t *testing.T) {
	posix_mq.ForceRemoveQueue("testMQ.5")
	msgSize := int(unsafe.Sizeof(TestMsg{}))
	sflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqtSend, err := createMQ("/testMQ.5", sflag, 10, msgSize)
	assert.Nil(t, err)
	rflag := posix_mq.O_RDONLY | posix_mq.O_CREAT
	mqtRecv, err := createMQ("/testMQ.5", rflag, 10, msgSize)
	assert.Nil(t, err)
	buf := bytes.NewBuffer(make([]byte, msgSize))
	for i := 1; i <= 10; i++ {
		msg := TestMsg{
			Type: uint8(i),
		}
		buf.Reset()
		err := binary.Write(buf, binary.LittleEndian, msg)
		assert.Nil(t, err)
		err = mqtSend.Send(buf.Bytes(), uint(i))
		assert.Nil(t, err)
		recvMsg, prio, err := mqtRecv.Receive()
		assert.Nil(t, err)
		assert.Equal(t, uint(i), prio)
		assert.Equal(t, uint8(i), recvMsg[0])
	}
	err = mqtSend.Unlink()
	assert.Nil(t, err)
	err = mqtRecv.Unlink()
	assert.NotNil(t, err)
	assert.Equal(t, syscall.ENOENT, err.(syscall.Errno))
}

// Send msg size larger than MsgSize send() wll return syscall.EMSGSIZE
func TestSendMsgTooLong(t *testing.T) {
	posix_mq.ForceRemoveQueue("testMQ.6")
	msgSize := int(unsafe.Sizeof(TestMsg{}))
	sflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqt, err := createMQ("/testMQ.6", sflag, 10, msgSize-1)
	assert.Nil(t, err)
	assert.NotNil(t, mqt)
	buf := bytes.NewBuffer(make([]byte, msgSize))
	err = binary.Write(buf, binary.LittleEndian, TestMsg{})
	assert.Nil(t, err)
	err = mqt.Send(buf.Bytes(), 0)
	assert.NotNil(t, err)
	assert.Equal(t, syscall.EMSGSIZE, err.(syscall.Errno))
	err = mqt.Unlink()
	assert.Nil(t, err)
}

// recv msg size smaller than MsgSize, receive() return syscall.EMSGSIZE
func TestRecvMsgTooShort(t *testing.T) {
	posix_mq.ForceRemoveQueue("testMQ.5")
	msgSize := int(unsafe.Sizeof(TestMsg{}))
	sflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqtSend, err := createMQ("/testMQ.5", sflag, 10, msgSize)
	assert.Nil(t, err)
	rflag := posix_mq.O_RDONLY | posix_mq.O_CREAT
	mqtRecv, err := createMQ("/testMQ.5", rflag, 10, msgSize-1)
	assert.Nil(t, err)
	buf := bytes.NewBuffer(make([]byte, msgSize))

	msg := TestMsg{
		Type: uint8(10),
	}
	buf.Reset()
	err = binary.Write(buf, binary.LittleEndian, msg)
	assert.Nil(t, err)
	err = mqtSend.Send(buf.Bytes(), 0)
	assert.Nil(t, err)
	recvMsg, prio, err := mqtRecv.Receive()
	assert.Equal(t, 0, len(recvMsg))
	assert.NotNil(t, err)
	assert.Equal(t, syscall.EMSGSIZE, err.(syscall.Errno))
	assert.Equal(t, uint(0), prio)

	err = mqtSend.Unlink()
	assert.Nil(t, err)
	err = mqtRecv.Unlink()
	assert.NotNil(t, err)
	assert.Equal(t, syscall.ENOENT, err.(syscall.Errno))
}

// with non-blocking queue, while queue full, send() will return syscall.EAGAIN
func TestSendwithNonblocking(t *testing.T) {
	posix_mq.ForceRemoveQueue("testMQ.5")
	msgSize := int(unsafe.Sizeof(TestMsg{}))
	sflag := posix_mq.O_WRONLY | posix_mq.O_CREAT | posix_mq.O_NONBLOCK
	mqtSend, err := createMQ("/testMQ.5", sflag, 1, msgSize)
	assert.Nil(t, err)

	buf := bytes.NewBuffer(make([]byte, msgSize))
	buf.Reset()
	err = binary.Write(buf, binary.LittleEndian, TestMsg{Type: uint8(1)})
	assert.Nil(t, err)
	err = mqtSend.Send(buf.Bytes(), 0)
	assert.Nil(t, err)
	buf.Reset()
	err = binary.Write(buf, binary.LittleEndian, TestMsg{Type: uint8(2)})
	assert.Nil(t, err)
	err = mqtSend.Send(buf.Bytes(), 1)
	assert.NotNil(t, err)
	assert.Equal(t, syscall.EAGAIN, err.(syscall.Errno))
	err = mqtSend.Unlink()
	assert.Nil(t, err)
}

// for empty non-blocking queue, receive() will return syscall.EAGAIN
func TestRecvEmptyQueueWithNonblocking(t *testing.T) {
	posix_mq.ForceRemoveQueue("testMQ.5")
	msgSize := int(unsafe.Sizeof(TestMsg{}))
	sflag := posix_mq.O_RDONLY | posix_mq.O_CREAT | posix_mq.O_NONBLOCK
	mqtRecv, err := createMQ("/testMQ.5", sflag, 1, msgSize)
	assert.Nil(t, err)

	msg, prio, err := mqtRecv.Receive()
	assert.NotNil(t, err)
	assert.Equal(t, uint(0), prio)
	assert.Equal(t, 0, len(msg))
	assert.EqualError(t, syscall.EAGAIN, err.Error())
}

func createMQ(name string, flag int, maxMsg int, msgSize int) (*posix_mq.MessageQueue, error) {
	attr := posix_mq.MessageQueueAttribute{
		MaxMsg:  maxMsg,
		MsgSize: msgSize,
	}
	return posix_mq.NewMessageQueue(name, flag, 666, &attr)
}

type TestMsg struct {
	Type   uint8
	Length uint8
	Data   [21]byte
}

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

func TestOpenMQWithOutCreatePermission(t *testing.T) {
	oflag := posix_mq.O_WRONLY
	mqt, err := posix_mq.NewMessageQueue("/testName", oflag, 666, nil)
	assert.Nil(t, mqt)
	mqErr, ok := err.(syscall.Errno)
	assert.True(t, ok)
	assert.Equal(t, syscall.ENOENT, mqErr)
}

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
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqt, err := createMQ("/testMQ.1", oflag, 19, 8196)
	assert.Nil(t, err)
	assert.NotNil(t, mqt)
	err = mqt.Unlink()
	assert.Nil(t, err)
}

func TestOpenExistMQ(t *testing.T) {
	oflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqt1, err := createMQ("/testMQ.1", oflag, 10, 8196)
	assert.Nil(t, err)
	assert.NotNil(t, mqt1)
	mqt2, err := createMQ("testMQ.2", oflag, 10, 8196)
	assert.Nil(t, err)
	assert.NotNil(t, mqt2)
	err = mqt1.Unlink()
	assert.Nil(t, err)

	err = mqt2.Unlink()
	assert.NotNil(t, err)
	assert.Equal(t, syscall.ENOENT, err.(syscall.Errno))
}

func TestExamQueueExist(t *testing.T) {
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
	msgSize := int(unsafe.Sizeof(TestMsg{}))
	var buf bytes.Buffer
	err := binary.Write(&buf, binary.LittleEndian, TestMsg{})
	assert.Nil(t, err)
	assert.Equal(t, msgSize, len(buf.Bytes()))

	sflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqtSend, err := createMQ("/testMQ.4", sflag, 10, msgSize)
	assert.Nil(t, err)
	assert.NotNil(t, mqtSend)

	for i := 1; i <= 10; i++ {
		err = mqtSend.Send(buf.Bytes(), 0)
		assert.Nil(t, err)
	}
	err = mqtSend.Unlink()
	assert.Nil(t, err)
}

func TestSendAndReceive(t *testing.T) {
	msgSize := int(unsafe.Sizeof(TestMsg{}))
	sflag := posix_mq.O_WRONLY | posix_mq.O_CREAT
	mqtSend, err := createMQ("/testMQ.5", sflag, 10, msgSize)
	assert.Nil(t, err)
	rflag := posix_mq.O_RDONLY
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
